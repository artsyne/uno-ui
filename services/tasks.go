package services

import (
	"encoding/json"
	"fmt"

	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/utils"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
)

// 定时任务设置
func RunTask() {
	gocron.Every(10).Seconds().Do(EmailTask)
	gocron.Start()
}

// 发送email任务
func EmailTask() {
	tasks, err := models.GetJobTasks(10, "email")
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Debug(fmt.Sprintf("获取%d条邮件任务", len(tasks)))

	//获取配置
	config := settings.GetConfig()

	//循环发送邮件
	for _, val := range tasks {
		var obj models.EmailTask
		err = json.Unmarshal([]byte(val.Body), &obj)
		if err != nil {
			logrus.Error(err)
			continue
		}

		//发送邮件
		err = utils.SendToMail(config.Email.SMTPUser, config.Email.SMTPPwd, config.Email.SMTPServer,
			obj.To, obj.Subject, obj.Body, "html")
		logrus.Debug(fmt.Sprintf("邮件发送成功 to=%s subject=%s ", obj.To, obj.Subject))

		if err != nil {
			logrus.Error(err)
			//任务重试次数加1
			_, _ = models.UpdateJobTaskTryNum(val.JobId)
			continue
		}

		//删除记录到历史表
		_, err := val.RemoveToHistory()
		if err != nil {
			logrus.Error(fmt.Sprintf("任务job_id=%d，删除到历史表失败", val.JobId))
		}
	}
}
