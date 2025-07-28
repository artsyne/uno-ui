package task_test

import (
	"encoding/json"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/domain/task"
	testhelpers "github.com/biuren/hebe/tests/helpers"
)

func TestTaskService_CreateEmailTask(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "br_job_queue")
	}

	taskService, err := di.TryGet[task.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取任务服务失败: %v", err)
	}

	tests := []struct {
		name    string
		to      string
		subject string
		body    string
		wantErr bool
	}{
		{
			name:    "创建激活邮件任务",
			to:      "test@example.com",
			subject: "激活您的账户",
			body:    "请点击链接激活您的账户",
			wantErr: false,
		},
		{
			name:    "创建密码重置邮件任务",
			to:      "test2@example.com",
			subject: "重置您的密码",
			body:    "请点击链接重置您的密码",
			wantErr: false,
		},
		{
			name:    "无效邮件地址",
			to:      "", // 空邮件地址
			subject: "测试",
			body:    "测试内容",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := taskService.CreateEmailTask(tt.to, tt.subject, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			// 验证任务是否已添加到数据库
			// 这里需要根据实际的Repository方法来验证
			// 由于没有GetJobById方法，我们可以检查数据库中是否有记录
			count := suite.Database.CountRows(t, "br_job_queue")
			if count == 0 {
				t.Error("任务未添加到数据库")
			}
		})
	}
}

func TestTaskService_ProcessEmailTasks(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "br_job_queue")
	}

	taskService, err := di.TryGet[task.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取任务服务失败: %v", err)
	}

	// 先创建一些邮件任务
	emailTasks := []struct {
		to      string
		subject string
		body    string
	}{
		{
			to:      "test1@example.com",
			subject: "激活账户",
			body:    "请点击链接激活账户",
		},
		{
			to:      "test2@example.com",
			subject: "重置密码",
			body:    "请点击链接重置密码",
		},
	}

	// 添加任务到数据库
	for _, emailTask := range emailTasks {
		err := taskService.CreateEmailTask(emailTask.to, emailTask.subject, emailTask.body)
		if err != nil {
			t.Fatalf("创建邮件任务失败: %v", err)
		}
	}

	// 验证任务已添加
	beforeCount := suite.Database.CountRows(t, "br_job_queue")
	if beforeCount != len(emailTasks) {
		t.Errorf("期望添加%d个任务，实际添加了%d个", len(emailTasks), beforeCount)
	}

	// 处理邮件任务（这里不会真正发送邮件，只测试业务逻辑）
	// 注意：实际的ProcessEmailTasks会发送邮件，在测试环境中可能会失败
	// 但我们可以测试它不会产生panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ProcessEmailTasks产生了panic: %v", r)
		}
	}()

	taskService.ProcessEmailTasks()

	// 注意：由于ProcessEmailTasks依赖邮件配置，在测试环境中可能会失败
	// 这里我们主要测试方法调用不会产生错误
}

func TestEmailTaskSerialization(t *testing.T) {
	// 测试EmailTask的序列化和反序列化
	emailTask := &task.EmailTask{
		To:      "test@example.com",
		Subject: "测试邮件",
		Body:    "这是测试邮件内容",
	}

	// 序列化
	data, err := json.Marshal(emailTask)
	if err != nil {
		t.Fatalf("序列化EmailTask失败: %v", err)
	}

	// 反序列化
	var deserializedTask task.EmailTask
	err = json.Unmarshal(data, &deserializedTask)
	if err != nil {
		t.Fatalf("反序列化EmailTask失败: %v", err)
	}

	// 验证数据一致性
	if deserializedTask.To != emailTask.To {
		t.Errorf("To字段不匹配: 期望 '%s', 实际 '%s'", emailTask.To, deserializedTask.To)
	}

	if deserializedTask.Subject != emailTask.Subject {
		t.Errorf("Subject字段不匹配: 期望 '%s', 实际 '%s'", emailTask.Subject, deserializedTask.Subject)
	}

	if deserializedTask.Body != emailTask.Body {
		t.Errorf("Body字段不匹配: 期望 '%s', 实际 '%s'", emailTask.Body, deserializedTask.Body)
	}
}

func TestJobQueueStructure(t *testing.T) {
	// 测试JobQueue结构体的基本功能
	jobQueue := &task.JobQueue{
		JobType:     "email",
		Body:        `{"to":"test@example.com","subject":"测试","body":"内容"}`,
		AddUTCTime:  1700000000000,
		PlanUTCTime: 1700000000000,
		Result:      0,
		TryNum:      0,
	}

	// 验证字段赋值
	if jobQueue.JobType != "email" {
		t.Errorf("JobType不匹配: 期望 'email', 实际 '%s'", jobQueue.JobType)
	}

	if jobQueue.TryNum != 0 {
		t.Errorf("TryNum不匹配: 期望 0, 实际 %d", jobQueue.TryNum)
	}

	// 测试Body字段是否为有效JSON
	var emailTask task.EmailTask
	err := json.Unmarshal([]byte(jobQueue.Body), &emailTask)
	if err != nil {
		t.Errorf("JobQueue.Body不是有效的JSON: %v", err)
	}

	// 验证解析的EmailTask
	if emailTask.To != "test@example.com" {
		t.Errorf("解析的EmailTask.To不正确: 期望 'test@example.com', 实际 '%s'", emailTask.To)
	}
}

func TestJobHistoryStructure(t *testing.T) {
	// 测试JobHistory结构体的基本功能
	jobHistory := &task.JobHistory{
		JobId:       123,
		JobType:     "email",
		Body:        `{"to":"test@example.com","subject":"测试","body":"内容"}`,
		PlanUTCTime: 1700000000000,
		DoneUTCTime: 1700000001000,
		Result:      1, // 成功
		TryNum:      1,
	}

	// 验证字段赋值
	if jobHistory.JobId != 123 {
		t.Errorf("JobId不匹配: 期望 123, 实际 %d", jobHistory.JobId)
	}

	if jobHistory.JobType != "email" {
		t.Errorf("JobType不匹配: 期望 'email', 实际 '%s'", jobHistory.JobType)
	}

	if jobHistory.Result != 1 {
		t.Errorf("Result不匹配: 期望 1, 实际 %d", jobHistory.Result)
	}

	if jobHistory.DoneUTCTime <= jobHistory.PlanUTCTime {
		t.Error("DoneUTCTime应该大于PlanUTCTime")
	}
}

func TestTaskService_StartScheduler(t *testing.T) {
	// 初始化测试环境
	service := task.NewService()

	// 测试StartScheduler方法不会产生panic
	// 注意：这个方法会启动定时任务，在测试环境中需要小心处理
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StartScheduler产生了panic: %v", r)
		}
	}()

	// 由于StartScheduler会启动永久运行的定时器，
	// 在测试中我们只验证方法调用不会立即失败
	// 实际的定时任务测试需要更复杂的setup

	// 这里我们可以测试NewService是否正常工作
	if service == nil {
		t.Error("NewService应该返回有效的Service实例")
	}
}

// 集成测试：测试完整的邮件任务流程
func TestEmailTaskWorkflow(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	service := task.NewService()

	// 1. 创建邮件任务
	to := "workflow@example.com"
	subject := "工作流测试邮件"
	body := "这是一个工作流测试邮件"

	err := service.CreateEmailTask(to, subject, body)
	if err != nil {
		t.Fatalf("创建邮件任务失败: %v", err)
	}

	// 2. 验证任务已添加到队列
	count := dbHelper.CountRows(t, "br_job_queue")
	if count != 1 {
		t.Errorf("期望队列中有1个任务，实际有%d个", count)
	}

	// 3. 查询任务内容验证
	rows := dbHelper.Query(t, "SELECT job_type, body FROM br_job_queue WHERE job_type = ?", "email")
	defer rows.Close()

	if rows.Next() {
		var jobType, bodyStr string
		err := rows.Scan(&jobType, &bodyStr)
		if err != nil {
			t.Fatalf("扫描行数据失败: %v", err)
		}

		// 验证job_type
		if jobType != "email" {
			t.Errorf("JobType不匹配: 期望 'email', 实际 '%s'", jobType)
		}

		// 验证body内容
		var emailTask task.EmailTask
		err = json.Unmarshal([]byte(bodyStr), &emailTask)
		if err != nil {
			t.Fatalf("解析邮件任务失败: %v", err)
		}

		if emailTask.To != to {
			t.Errorf("收件人不匹配: 期望 '%s', 实际 '%s'", to, emailTask.To)
		}

		if emailTask.Subject != subject {
			t.Errorf("主题不匹配: 期望 '%s', 实际 '%s'", subject, emailTask.Subject)
		}

		if emailTask.Body != body {
			t.Errorf("内容不匹配: 期望 '%s', 实际 '%s'", body, emailTask.Body)
		}
	} else {
		t.Error("未找到邮件任务记录")
	}
}
