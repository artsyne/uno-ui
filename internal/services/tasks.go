package services

import (
	"encoding/json"
	"fmt"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/modules/user"
	"github.com/biuren/hebe/pkg"

	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
)

// 定时任务设置
func RunTask() {
	logrus.Info("开始初始化定时任务...")

	// 设置定时任务
	gocron.Every(10).Seconds().Do(EmailTask)

	// 启动定时任务
	gocron.Start()

	logrus.Info("定时任务已启动，每10秒执行一次邮件任务")
}

// 发送email任务
func EmailTask() {
	// 创建用户仓库实例
	userRepo := user.NewRepository()

	tasks, err := userRepo.GetJobTasks(10, "email")
	if err != nil {
		logrus.Error("获取邮件任务失败:", err)
		return
	}
	logrus.Info(fmt.Sprintf("获取到%d条邮件任务", len(tasks)))

	//获取配置
	config := configs.GetConfig()

	//循环发送邮件
	for _, val := range tasks {
		var obj user.EmailTask
		err = json.Unmarshal([]byte(val.Body), &obj)
		if err != nil {
			logrus.Error("解析邮件任务失败:", err)
			continue
		}

		//发送邮件
		err = pkg.SendToMail(config.Email.SMTPUser, config.Email.SMTPPassword, config.Email.SMTPServer,
			obj.To, obj.Subject, obj.Body, "html")

		if err != nil {
			logrus.Error("邮件发送失败:", err)
			//任务重试次数加1
			_, _ = userRepo.UpdateJobTaskTryNum(val.JobId)
			continue
		}

		logrus.Info(fmt.Sprintf("邮件发送成功 to=%s subject=%s", obj.To, obj.Subject))

		//删除记录到历史表
		_, err := userRepo.RemoveJobToHistory(&val)
		if err != nil {
			logrus.Error(fmt.Sprintf("任务job_id=%d，删除到历史表失败", val.JobId))
		}
	}
}
