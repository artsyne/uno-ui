package task

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/jasonlvhit/gocron"
)

// Service 任务业务逻辑层实现
// 负责处理所有与任务相关的业务逻辑和规则验证
type Service struct {
	config *config.Config       // 配置管理器
	repo   RepositoryInterface  // 任务数据访问层接口
}

// NewService 创建新的任务业务逻辑层实例
// 参数:
//   - config: 配置管理器
//   - repo: 任务数据访问层接口
// 返回:
//   - *Service: 任务业务逻辑层实例
func NewService(config *config.Config, repo RepositoryInterface) *Service {
	return &Service{
		config: config,
		repo:   repo,
	}
}

// ===================== 接口实现方法 =====================

// StartScheduler 启动定时任务调度器
// 初始化并启动所有定时任务
// 注意: 该方法会清除之前的任务，避免重复注册
func (s *Service) StartScheduler() {
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
// 该方法会被定时调度器定期调用
func (s *Service) ProcessEmailTasks() {
	tasks, err := s.repo.GetJobTasks(10, "email")
	if err != nil {
		slog.Error("获取邮件任务失败", "error", err)
		return
	}

	// 没有获取任务退出
	tasksTotals := len(tasks)
	slog.Debug("获取邮件任务", "count", tasksTotals)
	if tasksTotals == 0 {
		return
	}
	
	//循环发送邮件
	for _, val := range tasks {
		s.processEmailTask(val, s.config)
	}
}

// CreateEmailTask 创建邮件任务
// 将邮件任务添加到队列中等待处理
// 参数:
//   - to: 收件人邮箱地址
//   - subject: 邮件主题
//   - body: 邮件正文内容
// 返回:
//   - error: 创建失败时返回错误
func (s *Service) CreateEmailTask(to, subject, body string) error {
	emailTask := &EmailTask{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	taskBody, err := json.Marshal(emailTask)
	if err != nil {
		return errors.ValidationError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "email_serialization_failed"})
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
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "email_job_creation_failed"})
	}
	return nil
}

// ===================== 私有辅助方法 =====================

// processEmailTask 处理单个邮件任务
// 解析任务内容，发送邮件，处理结果
// 参数:
//   - task: 待处理的邮件任务
//   - config: 配置管理器
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

	slog.Debug("邮件发送成功", "job_id", task.JobId, "to", emailTask.To, "subject", emailTask.Subject)

	//删除记录到历史表
	_, err = s.repo.RemoveJobToHistory(&task)
	if err != nil {
		slog.Error("任务删除到历史表失败", "job_id", task.JobId, "error", err)
	}
}