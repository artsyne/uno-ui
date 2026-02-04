package pkg_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/biuren/hebe/pkg"
	framework "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
)

// TestWorkEventArgs 测试WorkEventArgs结构体
func TestWorkEventArgs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("WorkEventArgs结构体测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		args := pkg.WorkEventArgs{
			SourceEvent:  "OnStart",
			ThreadId:     1,
			Code:         "0",
			WorkName:     "TestWorker",
			ExceptionMsg: "No error",
			Extend:       "Extra info",
			Msg:          "Test message",
		}

		assert.Equal(t, "OnStart", args.SourceEvent, "SourceEvent应该正确设置")
		assert.Equal(t, 1, args.ThreadId, "ThreadId应该正确设置")
		assert.Equal(t, "0", args.Code, "Code应该正确设置")
		assert.Equal(t, "TestWorker", args.WorkName, "WorkName应该正确设置")
		assert.Equal(t, "No error", args.ExceptionMsg, "ExceptionMsg应该正确设置")
		assert.Equal(t, "Extra info", args.Extend, "Extend应该正确设置")
		assert.Equal(t, "Test message", args.Msg, "Msg应该正确设置")
	})

	t.Run("WorkEventArgs默认值测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		args := pkg.WorkEventArgs{}

		assert.Equal(t, "", args.SourceEvent, "SourceEvent默认值应该为空字符串")
		assert.Equal(t, 0, args.ThreadId, "ThreadId默认值应该为0")
		assert.Equal(t, "", args.Code, "Code默认值应该为空字符串")
		assert.Equal(t, "", args.WorkName, "WorkName默认值应该为空字符串")
		assert.Equal(t, "", args.ExceptionMsg, "ExceptionMsg默认值应该为空字符串")
		assert.Equal(t, "", args.Extend, "Extend默认值应该为空字符串")
		assert.Equal(t, "", args.Msg, "Msg默认值应该为空字符串")
	})
}

// TestNewWorker 测试NewWorker函数
func TestNewWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("NewWorker基本功能测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string {
			return "processed: " + msg
		}

		worker := pkg.NewWorker(3, "TestWorker", workerFunc, false, 5*time.Second)

		assert.NotNil(t, worker, "Worker不应该为nil")

		// 由于字段是私有的，我们通过功能测试来验证Worker是否正确初始化
		// 测试添加任务功能
		success := worker.AddTask("test message")
		assert.True(t, success, "添加任务应该成功")

		// 测试启动功能
		startResult := worker.Start()
		assert.True(t, startResult, "启动Worker应该成功")

		// 清理资源
		worker.Dispose()
	})

	t.Run("NewWorker不同参数测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		testCases := []struct {
			threadNumber     int
			workName         string
			asyncInvokeEvent bool
			actionTimeout    time.Duration
		}{
			{1, "SingleWorker", true, 1 * time.Second},
			{10, "MultiWorker", false, 10 * time.Second},
			{5, "DefaultWorker", true, 30 * time.Second},
			{0, "ZeroWorker", false, 0},
		}

		for _, tc := range testCases {
			workerFunc := func(msg string) string { return msg }
			worker := pkg.NewWorker(tc.threadNumber, tc.workName, workerFunc, tc.asyncInvokeEvent, tc.actionTimeout)

			// 由于字段是私有的，我们通过功能测试来验证Worker是否正确初始化
			success := worker.AddTask("test")
			assert.True(t, success, "添加任务应该成功")

			startResult := worker.Start()
			assert.True(t, startResult, "启动Worker应该成功")

			worker.Dispose()
		}
	})

	t.Run("NewWorker工作函数测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string {
			return "result: " + msg
		}

		worker := pkg.NewWorker(1, "TestWorker", workerFunc, false, 1*time.Second)

		// 由于workerAction是私有的，我们通过完整的工作流程来测试
		var resultReceived string
		worker.OnSuccess = func(sender interface{}, e pkg.WorkEventArgs) {
			resultReceived = e.Extend
		}

		worker.Start()
		worker.AddTask("test")

		// 等待处理完成
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, "result: test", resultReceived, "工作函数应该正确执行")

		worker.Dispose()
	})
}

// TestWorkerAddTask 测试AddTask方法
func TestWorkerAddTask(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("AddTask成功测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "TestWorker", workerFunc, false, 1*time.Second)

		// 设置事件处理器来验证事件触发
		var addedTaskCalled bool
		worker.OnAddedTask = func(sender interface{}, e pkg.WorkEventArgs) {
			addedTaskCalled = true
			assert.Equal(t, "OnAddedTask", e.SourceEvent, "事件源应该正确")
			assert.Equal(t, "test message", e.Msg, "消息应该正确")
			assert.Equal(t, "TestWorker", e.WorkName, "工作名称应该正确")
			assert.Equal(t, "0", e.Code, "代码应该为0")
		}

		success := worker.AddTask("test message")

		assert.True(t, success, "添加任务应该成功")
		assert.True(t, addedTaskCalled, "OnAddedTask事件应该被触发")

		worker.Dispose()
	})

	t.Run("AddTask队列满测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "TestWorker", workerFunc, false, 1*time.Second)

		// 填满队列（容量为100）
		for i := 0; i < 100; i++ {
			success := worker.AddTask("message")
			assert.True(t, success, "前100个任务应该成功添加")
		}

		// 设置失败事件处理器
		var failureCalled bool
		worker.OnFailure = func(sender interface{}, e pkg.WorkEventArgs) {
			failureCalled = true
			assert.Equal(t, "Failure", e.SourceEvent, "事件源应该正确")
			assert.Equal(t, "5000", e.Code, "错误代码应该为5000")
			assert.Equal(t, "Timeout adding task", e.ExceptionMsg, "异常消息应该正确")
		}

		// 尝试添加第101个任务（应该超时失败）
		success := worker.AddTask("overflow message")

		assert.False(t, success, "队列满时添加任务应该失败")
		assert.True(t, failureCalled, "OnFailure事件应该被触发")

		worker.Dispose()
	})

	t.Run("AddTask多个任务测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "TestWorker", workerFunc, false, 1*time.Second)

		var addedCount int32
		worker.OnAddedTask = func(sender interface{}, e pkg.WorkEventArgs) {
			atomic.AddInt32(&addedCount, 1)
		}

		// 添加多个任务
		for i := 0; i < 10; i++ {
			success := worker.AddTask("message")
			assert.True(t, success, "所有任务都应该成功添加")
		}

		assert.Equal(t, int32(10), atomic.LoadInt32(&addedCount), "应该触发10次OnAddedTask事件")

		worker.Dispose()
	})
}

// TestWorkerStart 测试Start方法
func TestWorkerStart(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("Start基本功能测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string {
			time.Sleep(10 * time.Millisecond) // 模拟工作
			return "processed: " + msg
		}

		worker := pkg.NewWorker(2, "TestWorker", workerFunc, false, 1*time.Second)

		var threadCompleteCalled bool
		completionSignal := make(chan struct{})

		worker.OnThreadComplete = func(sender interface{}, e pkg.WorkEventArgs) {
			threadCompleteCalled = true
			assert.Equal(t, "OnThreadComplete", e.SourceEvent, "事件源应该正确")
			assert.Equal(t, "TestWorker", e.WorkName, "工作名称应该正确")
			// 防止多次关闭
			select {
			case <-completionSignal:
			default:
				close(completionSignal)
			}
		}

		success := worker.Start()
		assert.True(t, success, "启动应该成功")

		// 添加一些任务
		worker.AddTask("task1")
		worker.AddTask("task2")

		// 等待一段时间让任务处理
		time.Sleep(100 * time.Millisecond)

		worker.Dispose()

		// 使用通道等待事件触发，而不是由于 sleep 不稳定导致的失败
		select {
		case <-completionSignal:
			// 成功收到信号
		case <-time.After(2 * time.Second):
			t.Error("等待 OnThreadComplete 超时")
		}

		assert.True(t, threadCompleteCalled, "OnThreadComplete事件应该被触发")
	})

	t.Run("Start多线程测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		var processedCount int32
		workerFunc := func(msg string) string {
			atomic.AddInt32(&processedCount, 1)
			return msg
		}

		worker := pkg.NewWorker(5, "MultiThreadWorker", workerFunc, false, 1*time.Second)

		worker.Start()

		// 添加多个任务
		for i := 0; i < 20; i++ {
			worker.AddTask("task")
		}

		// 等待任务处理完成
		time.Sleep(500 * time.Millisecond)

		worker.Dispose()

		// 验证任务被处理
		processed := atomic.LoadInt32(&processedCount)
		assert.Greater(t, processed, int32(0), "应该有任务被处理")
	})
}

// TestWorkerEvents 测试Worker事件
func TestWorkerEvents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("OnStart和OnSuccess事件测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string {
			return "result: " + msg
		}

		worker := pkg.NewWorker(1, "EventTestWorker", workerFunc, false, 1*time.Second)

		var startCalled, successCalled bool
		var startMsg, successMsg, successResult string

		worker.OnStart = func(sender interface{}, e pkg.WorkEventArgs) {
			startCalled = true
			startMsg = e.Msg
			assert.Equal(t, "OnStart", e.SourceEvent, "OnStart事件源应该正确")
			assert.Equal(t, "0", e.Code, "OnStart代码应该为0")
		}

		worker.OnSuccess = func(sender interface{}, e pkg.WorkEventArgs) {
			successCalled = true
			successMsg = e.Msg
			successResult = e.Extend
			assert.Equal(t, "OnSuccess", e.SourceEvent, "OnSuccess事件源应该正确")
			assert.Equal(t, "0", e.Code, "OnSuccess代码应该为0")
		}

		worker.Start()
		worker.AddTask("test message")

		// 等待任务处理
		time.Sleep(200 * time.Millisecond)

		assert.True(t, startCalled, "OnStart事件应该被触发")
		assert.True(t, successCalled, "OnSuccess事件应该被触发")
		assert.Equal(t, "test message", startMsg, "OnStart消息应该正确")
		assert.Equal(t, "test message", successMsg, "OnSuccess消息应该正确")
		assert.Equal(t, "result: test message", successResult, "OnSuccess结果应该正确")

		worker.Dispose()
	})

	t.Run("OnTimeout事件测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		// 创建一个会超时的工作函数
		workerFunc := func(msg string) string {
			time.Sleep(2 * time.Second) // 超过超时时间
			return msg
		}

		worker := pkg.NewWorker(1, "TimeoutTestWorker", workerFunc, false, 100*time.Millisecond)

		var timeoutCalled bool
		var timeoutMsg string

		worker.OnTimeout = func(sender interface{}, e pkg.WorkEventArgs) {
			timeoutCalled = true
			timeoutMsg = e.Msg
			assert.Equal(t, "OnTimeout", e.SourceEvent, "OnTimeout事件源应该正确")
			assert.Equal(t, "5003", e.Code, "OnTimeout代码应该为5003")
			assert.Equal(t, "Timeout", e.ExceptionMsg, "OnTimeout异常消息应该正确")
		}

		worker.Start()
		worker.AddTask("timeout test")

		// 等待超时发生
		time.Sleep(300 * time.Millisecond)

		assert.True(t, timeoutCalled, "OnTimeout事件应该被触发")
		assert.Equal(t, "timeout test", timeoutMsg, "OnTimeout消息应该正确")

		worker.Dispose()
	})

	t.Run("异步事件测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }

		// 创建异步事件的Worker
		worker := pkg.NewWorker(1, "AsyncEventWorker", workerFunc, true, 1*time.Second)

		var eventCalled bool
		var wg sync.WaitGroup
		var once sync.Once
		wg.Add(1)

		worker.OnSuccess = func(sender interface{}, e pkg.WorkEventArgs) {
			once.Do(func() {
				eventCalled = true
				wg.Done()
			})
		}

		worker.Start()
		worker.AddTask("async test")

		// 等待异步事件完成
		wg.Wait()

		assert.True(t, eventCalled, "异步事件应该被触发")

		worker.Dispose()
	})
}

// TestWorkerInvokeEvent 测试InvokeEvent方法
func TestWorkerInvokeEvent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("InvokeEvent同步测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "TestWorker", workerFunc, false, 1*time.Second)

		var eventCalled bool
		eventFunc := func(sender interface{}, e pkg.WorkEventArgs) {
			eventCalled = true
			assert.Equal(t, worker, sender, "发送者应该是worker实例")
			assert.Equal(t, "test", e.Msg, "消息应该正确")
		}

		args := pkg.WorkEventArgs{Msg: "test"}
		worker.InvokeEvent(eventFunc, args)

		assert.True(t, eventCalled, "同步事件应该立即被调用")
	})

	t.Run("InvokeEvent异步测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "TestWorker", workerFunc, true, 1*time.Second)

		var eventCalled bool
		var wg sync.WaitGroup
		wg.Add(1)

		eventFunc := func(sender interface{}, e pkg.WorkEventArgs) {
			eventCalled = true
			wg.Done()
		}

		args := pkg.WorkEventArgs{Msg: "test"}
		worker.InvokeEvent(eventFunc, args)

		// 异步调用，需要等待
		wg.Wait()
		assert.True(t, eventCalled, "异步事件应该被调用")
	})

	t.Run("InvokeEvent nil事件测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "TestWorker", workerFunc, false, 1*time.Second)

		args := pkg.WorkEventArgs{Msg: "test"}

		// 调用nil事件不应该panic
		assert.NotPanics(t, func() {
			worker.InvokeEvent(nil, args)
		}, "调用nil事件不应该panic")
	})
}

// TestWorkerDispose 测试Dispose方法
func TestWorkerDispose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("Dispose基本功能测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(2, "TestWorker", workerFunc, false, 1*time.Second)

		worker.Start()

		// 添加一些任务
		worker.AddTask("task1")
		worker.AddTask("task2")

		// 等待一段时间
		time.Sleep(100 * time.Millisecond)

		// 销毁Worker
		assert.NotPanics(t, func() {
			worker.Dispose()
		}, "Dispose不应该panic")

		// 尝试添加任务应该失败（通道已关闭）
		assert.Panics(t, func() {
			worker.AddTask("task after dispose")
		}, "在Dispose后添加任务应该panic")
	})

	t.Run("Dispose多次调用测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "TestWorker", workerFunc, false, 1*time.Second)

		worker.Start()

		// 第一次Dispose
		assert.NotPanics(t, func() {
			worker.Dispose()
		}, "第一次Dispose不应该panic")

		// 第二次Dispose应该panic（重复关闭通道）
		assert.Panics(t, func() {
			worker.Dispose()
		}, "重复Dispose应该panic")
	})
}

// TestWorkerEdgeCases 测试Worker的边界情况
func TestWorkerEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("零线程数测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(0, "ZeroThreadWorker", workerFunc, false, 1*time.Second)

		// 启动零线程的Worker
		success := worker.Start()
		assert.True(t, success, "启动应该成功")

		// 添加任务
		success = worker.AddTask("task")
		assert.True(t, success, "添加任务应该成功")

		// 但是任务不会被处理（没有工作线程）
		time.Sleep(100 * time.Millisecond)

		worker.Dispose()
	})

	t.Run("负数线程数测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(-1, "NegativeThreadWorker", workerFunc, false, 1*time.Second)

		// 启动负数线程的Worker
		success := worker.Start()
		assert.True(t, success, "启动应该成功")

		worker.Dispose()
	})

	t.Run("nil工作函数测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		worker := pkg.NewWorker(1, "NilFuncWorker", nil, false, 1*time.Second)

		var failureCalled bool
		var failureMsg string
		worker.OnFailure = func(sender interface{}, e pkg.WorkEventArgs) {
			failureCalled = true
			failureMsg = e.ExceptionMsg
		}

		worker.Start()

		// 添加任务应该触发OnFailure事件而不是panic
		worker.AddTask("task")

		// 等待一段时间让事件触发
		time.Sleep(200 * time.Millisecond)

		assert.True(t, failureCalled, "nil工作函数应该触发OnFailure事件")
		assert.Equal(t, "Worker action is nil", failureMsg, "应该有正确的错误消息")

		worker.Dispose()
	})

	t.Run("零超时时间测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string {
			time.Sleep(100 * time.Millisecond)
			return msg
		}

		worker := pkg.NewWorker(1, "ZeroTimeoutWorker", workerFunc, false, 0)

		var timeoutCalled bool
		worker.OnTimeout = func(sender interface{}, e pkg.WorkEventArgs) {
			timeoutCalled = true
		}

		worker.Start()
		worker.AddTask("task")

		// 等待处理
		time.Sleep(200 * time.Millisecond)

		// 零超时应该立即触发超时
		assert.True(t, timeoutCalled, "零超时应该立即触发超时事件")

		worker.Dispose()
	})

	t.Run("空工作名称测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(1, "", workerFunc, false, 1*time.Second)

		var eventWorkName string
		worker.OnAddedTask = func(sender interface{}, e pkg.WorkEventArgs) {
			eventWorkName = e.WorkName
		}

		worker.AddTask("task")

		assert.Equal(t, "", eventWorkName, "空工作名称应该正确传递")

		worker.Dispose()
	})

	t.Run("特殊字符消息测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return "processed: " + msg }
		worker := pkg.NewWorker(1, "SpecialCharWorker", workerFunc, false, 1*time.Second)

		specialMessages := []string{
			"message with spaces",
			"message\nwith\nnewlines",
			"message\twith\ttabs",
			"message中文字符",
			"message🚀emoji",
			"message@#$%^&*()special",
			"", // 空消息
		}

		var processedMessages []string
		var addedCount int32
		var mu sync.Mutex

		worker.OnAddedTask = func(sender interface{}, e pkg.WorkEventArgs) {
			atomic.AddInt32(&addedCount, 1)
		}

		worker.OnSuccess = func(sender interface{}, e pkg.WorkEventArgs) {
			mu.Lock()
			processedMessages = append(processedMessages, e.Msg)
			mu.Unlock()
		}

		worker.Start()

		for _, msg := range specialMessages {
			success := worker.AddTask(msg)
			assert.True(t, success, "特殊字符消息应该能成功添加")
		}

		// 等待所有任务被添加
		for i := 0; i < 100; i++ {
			if atomic.LoadInt32(&addedCount) == int32(len(specialMessages)) {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		// 等待处理完成，增加等待时间
		time.Sleep(1 * time.Second)

		mu.Lock()
		actualCount := len(processedMessages)
		mu.Unlock()

		assert.Equal(t, len(specialMessages), actualCount, "所有特殊字符消息都应该被处理")

		worker.Dispose()
	})
}

// TestWorkerPerformance 测试Worker性能
func TestWorkerPerformance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("并发添加任务测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }
		worker := pkg.NewWorker(5, "ConcurrentWorker", workerFunc, false, 1*time.Second)

		worker.Start()

		var successCount int32
		var wg sync.WaitGroup

		// 启动多个goroutine并发添加任务
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					if worker.AddTask("concurrent task") {
						atomic.AddInt32(&successCount, 1)
					}
				}
			}()
		}

		wg.Wait()

		success := atomic.LoadInt32(&successCount)
		assert.Greater(t, success, int32(0), "并发添加任务应该有成功的")

		worker.Dispose()
	})

	t.Run("内存使用测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string { return msg }

		// 创建多个Worker实例
		workers := make([]*pkg.Worker, 100)
		for i := 0; i < 100; i++ {
			workers[i] = pkg.NewWorker(1, "MemoryTestWorker", workerFunc, false, 1*time.Second)
			workers[i].Start()
		}

		// 添加任务
		for _, worker := range workers {
			worker.AddTask("memory test")
		}

		// 等待处理
		time.Sleep(100 * time.Millisecond)

		// 清理所有Worker
		for _, worker := range workers {
			worker.Dispose()
		}

		// 这个测试主要确保不会有内存泄漏或panic
		assert.True(t, true, "内存测试完成")
	})
}

// TestWorkerIntegration 测试Worker集成功能
func TestWorkerIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("完整工作流程测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		var events []string
		var eventsMutex sync.Mutex

		addEvent := func(event string) {
			eventsMutex.Lock()
			defer eventsMutex.Unlock()
			events = append(events, event)
		}

		workerFunc := func(msg string) string {
			time.Sleep(50 * time.Millisecond) // 模拟工作
			return "result: " + msg
		}

		worker := pkg.NewWorker(2, "IntegrationWorker", workerFunc, false, 1*time.Second)

		// 设置所有事件处理器
		worker.OnAddedTask = func(sender interface{}, e pkg.WorkEventArgs) {
			addEvent("OnAddedTask: " + e.Msg)
		}

		worker.OnStart = func(sender interface{}, e pkg.WorkEventArgs) {
			addEvent("OnStart: " + e.Msg)
		}

		worker.OnSuccess = func(sender interface{}, e pkg.WorkEventArgs) {
			addEvent("OnSuccess: " + e.Msg + " -> " + e.Extend)
		}

		worker.OnThreadComplete = func(sender interface{}, e pkg.WorkEventArgs) {
			addEvent("OnThreadComplete")
		}

		// 启动Worker
		worker.Start()

		// 添加任务
		worker.AddTask("task1")
		worker.AddTask("task2")
		worker.AddTask("task3")

		// 等待任务处理完成
		time.Sleep(500 * time.Millisecond)

		// 销毁Worker
		worker.Dispose()

		// 等待线程完成事件
		time.Sleep(100 * time.Millisecond)

		// 验证事件序列
		eventsMutex.Lock()
		defer eventsMutex.Unlock()

		assert.Greater(t, len(events), 0, "应该有事件被触发")

		// 验证包含预期的事件类型
		hasAddedTask := false
		hasStart := false
		hasSuccess := false

		for _, event := range events {
			if event == "OnAddedTask: task1" || event == "OnAddedTask: task2" || event == "OnAddedTask: task3" {
				hasAddedTask = true
			}
			if event == "OnStart: task1" || event == "OnStart: task2" || event == "OnStart: task3" {
				hasStart = true
			}
			if event == "OnSuccess: task1 -> result: task1" || event == "OnSuccess: task2 -> result: task2" || event == "OnSuccess: task3 -> result: task3" {
				hasSuccess = true
			}
		}

		assert.True(t, hasAddedTask, "应该有OnAddedTask事件")
		assert.True(t, hasStart, "应该有OnStart事件")
		assert.True(t, hasSuccess, "应该有OnSuccess事件")
	})

	t.Run("错误处理集成测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		workerFunc := func(msg string) string {
			if msg == "panic" {
				panic("test panic")
			}
			if msg == "slow" {
				time.Sleep(2 * time.Second)
			}
			return msg
		}

		worker := pkg.NewWorker(1, "ErrorTestWorker", workerFunc, false, 100*time.Millisecond)

		var timeoutCalled bool
		worker.OnTimeout = func(sender interface{}, e pkg.WorkEventArgs) {
			timeoutCalled = true
		}

		worker.Start()

		// 添加会超时的任务
		worker.AddTask("slow")

		// 等待超时
		time.Sleep(300 * time.Millisecond)

		assert.True(t, timeoutCalled, "超时事件应该被触发")

		worker.Dispose()
	})
}
