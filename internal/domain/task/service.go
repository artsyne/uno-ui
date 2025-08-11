package task

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/biuren/hebe/internal/utils/config"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
	"github.com/jasonlvhit/gocron"
)



// Service 任务业务逻辑层实现
// 负责处理所有与任务相关的业务逻辑和规则验证
type Service struct {
	repo   RepositoryInterface
	config *config.Config // 缓存配置，避免重复调用全局变量
}

// NewService 创建新的任务业务逻辑层实例
func NewService(repo RepositoryInterface, cfg *config.Config) *Service {

	return &Service{
		repo:   repo,
		config: cfg, // 缓存配置
	}
}

// ===================== 接口实现方法 =====================

// StartScheduler 启动定时任务调度器
// 初始化并启动所有定时任务
func (s *Service) StartScheduler() {
	slog.Info("开始初始化定时任务...")

	// 清除之前的任务，避免重复注册
	gocron.Clear()
	
	// 设置定时任务
	gocron.Every(10).Seconds().Do(s.ProcessEmailTasks)

	// 在后台启动定时任务，避免阻塞主线程
	go func() {
		gocron.Start()
	}()

	slog.Info("定时任务已启动，每10秒执行一次邮件任务")
}

// ProcessEmailTasks 处理邮件任务
// 获取待处理的邮件任务，发送邮件，处理结果
func (s *Service) ProcessEmailTasks() {
	tasks, err := s.repo.GetJobTasks(10, "email")
	if err != nil {
		slog.Error("获取邮件任务失败", "error", err)
		return
	}
	slog.Info("获取邮件任务", "count", len(tasks))

	//循环发送邮件
	for _, val := range tasks {
		s.processEmailTask(val, s.config)
	}
}

// CreateEmailTask 创建邮件任务
// 将邮件任务添加到队列中等待处理
func (s *Service) CreateEmailTask(to, subject, body string) error {
	emailTask := &EmailTask{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	taskBody, err := json.Marshal(emailTask)
	if err != nil {
		return errors.ValidationError.WithDetails(fmt.Sprintf("邮件任务序列化失败: %v", err))
	}

	nowTime := time.Now().Unix()
	jobQueue := &JobQueue{
		JobType:     "email",
		Body:        string(taskBody),
		Result:      0,
		TryNum:      0,
		AddUTCTime:  nowTime,
		PlanUTCTime: nowTime,
	}

	_, err = s.repo.InsertJobQueue(jobQueue)
	if err != nil {
		return errors.DatabaseError.WithDetails(fmt.Sprintf("邮件任务创建失败: %v", err))
	}
	return nil
}

// ===================== 私有辅助方法 =====================

// processEmailTask 处理单个邮件任务
// 解析任务内容，发送邮件，处理结果
func (s *Service) processEmailTask(task JobQueue, config *config.Config) {
	var emailTask EmailTask
	err := json.Unmarshal([]byte(task.Body), &emailTask)
	if err != nil {
		slog.Error("解析邮件任务失败", "error", err, "job_id", task.JobId)
		return
	}

	//发送邮件
	err = pkg.SendToMail(
		config.Email.SMTPUser,
		config.Email.SMTPPassword,
		config.Email.SMTPServer,
		emailTask.To,
		emailTask.Subject,
		emailTask.Body,
		"html",
	)

	if err != nil {
		slog.Error("邮件发送失败", "error", err, "job_id", task.JobId, "to", emailTask.To)
		//任务重试次数加1
		_, _ = s.repo.UpdateJobTaskTryNum(task.JobId)
		return
	}

	slog.Info("邮件发送成功", "job_id", task.JobId, "to", emailTask.To, "subject", emailTask.Subject)

	//删除记录到历史表
	_, err = s.repo.RemoveJobToHistory(&task)
	if err != nil {
		slog.Error("任务删除到历史表失败", "job_id", task.JobId, "error", err)
	}
}