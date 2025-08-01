package task

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/biuren/hebe/internal/utils/config"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
	"github.com/jasonlvhit/gocron"
)

// ServiceInterface 任务业务逻辑层接口
// 定义所有与任务相关的业务逻辑操作
type ServiceInterface interface {
	// StartScheduler 启动定时任务调度器
	StartScheduler()

	// ProcessEmailTasks 处理邮件任务
	ProcessEmailTasks()

	// CreateEmailTask 创建邮件任务
	CreateEmailTask(to, subject, body string) error
}

// service 任务业务逻辑层实现
// 负责处理所有与任务相关的业务逻辑和规则验证
type service struct {
	repo   RepositoryInterface
	config *config.Config // 缓存配置，避免重复调用全局变量
}

// NewService 创建新的任务业务逻辑层实例
func NewService() ServiceInterface {
	config := config.GetConfig()
	if config == nil {
		// 在初始化时就检查配置，如果配置未正确加载则panic
		// 这是合理的，因为没有配置应用无法正常运行
		panic("Failed to get config during task service initialization - config is nil")
	}

	return &service{
		repo:   NewRepository(),
		config: config, // 缓存配置
	}
}

// ===================== 接口实现方法 =====================

// StartScheduler 启动定时任务调度器
// 初始化并启动所有定时任务
func (s *service) StartScheduler() {
	slog.Info("开始初始化定时任务...")

	// 设置定时任务
	gocron.Every(10).Seconds().Do(s.ProcessEmailTasks)

	// 启动定时任务
	gocron.Start()

	slog.Info("定时任务已启动，每10秒执行一次邮件任务")
}

// ProcessEmailTasks 处理邮件任务
// 获取待处理的邮件任务，发送邮件，处理结果
func (s *service) ProcessEmailTasks() {
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
func (s *service) CreateEmailTask(to, subject, body string) error {
	emailTask := &EmailTask{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	taskBody, err := json.Marshal(emailTask)
	if err != nil {
		return errors.ValidationError.WithDetails(fmt.Sprintf("邮件任务序列化失败: %v", err))
	}

	nowTime := pkg.GetNowTimeUnixMilli()
	jobQueue := &JobQueue{
		JobType:     "email",
		Body:        string(taskBody),
		AddUTCTime:  nowTime,
		PlanUTCTime: nowTime,
		Result:      0,
		TryNum:      0,
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
func (s *service) processEmailTask(task JobQueue, config *config.Config) {
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