package pkg

import (
    "sync"
    "time"
)

// WorkEventArgs 用于传递事件的参数
type WorkEventArgs struct {
    SourceEvent  string // 源事件名称，例如 "OnStart" 或 "OnSuccess"
    ThreadId     int    // 线程ID，用于标识哪个线程触发了事件
    Code         string // 错误码，用于标识事件的状态或错误类型
    WorkName     string // 工作类名，用于标识当前的工作类
    ExceptionMsg string // 异常消息说明，用于描述错误的详细信息
    Extend       string // 扩展信息，可以用于传递额外的上下文信息
    Msg          string // 当前处理的消息，表示正在处理的任务内容
}

// Worker 结构体定义了一个多线程工作者容器
type Worker struct {
    threadNumber    int                      // 工作线程数，表示要启动的并发工作线程的数量
    workName        string                   // 工作类名称，用于标识当前的工作类
    workQueue       chan string              // 消息传递的队列，用于存储待处理的任务
    cts             chan struct{}            // 取消队列时使用的通道，用于通知工作线程停止工作
    workerAction    func(string) string      // 工作事件，定义了处理任务的函数
    asyncInvokeEvent bool                    // 是否异步执行事件，决定事件处理是同步还是异步
    actionTimeout   time.Duration            // 执行委托方法超时时间，表示任务处理的超时时间
    invertalTakeQueue time.Duration          // 设置获取队列频率，表示工作线程获取任务的间隔时间

    // 事件定义
    OnThreadComplete func(sender interface{}, e WorkEventArgs) // 当线程启动完毕时触发的事件
    OnAddedTask      func(sender interface{}, e WorkEventArgs) // 当成功提交一个任务执行时触发的事件
    OnStart          func(sender interface{}, e WorkEventArgs) // 方法开始时执行的事件
    OnSuccess        func(sender interface{}, e WorkEventArgs) // 执行成功时触发的事件
    OnFailure        func(sender interface{}, e WorkEventArgs) // 当任务失败时触发的事件
    OnTimeout        func(sender interface{}, e WorkEventArgs) // 当任务超时时触发的事件
}

// NewWorker 初始化 Worker
func NewWorker(threadNumber int, workName string, worker func(string) string, asyncInvokeEvent bool, actionTimeout time.Duration) *Worker {
    return &Worker{
        threadNumber:    threadNumber,          // 初始化工作线程数
        workName:        workName,              // 初始化工作类名称
        workQueue:       make(chan string, 100), // 初始化消息队列，容量为100
        cts:             make(chan struct{}),   // 初始化取消信号通道
        workerAction:    worker,                // 初始化工作事件处理函数
        asyncInvokeEvent: asyncInvokeEvent,     // 初始化是否异步执行事件
        actionTimeout:   actionTimeout,         // 初始化任务处理超时时间
        invertalTakeQueue: 100 * time.Millisecond, // 初始化获取队列频率，默认100毫秒
    }
}

// AddTask 提交一个消息到工作队列
func (w *Worker) AddTask(msg string) bool {
    select {
    case w.workQueue <- msg:
        // 触发 OnAddedTask 事件
        w.InvokeEvent(w.OnAddedTask, WorkEventArgs{
            Code:     "0",
            ThreadId: 0,
            WorkName: w.workName,
            Msg:      msg,
            SourceEvent: "OnAddedTask",
        })
        return true
    case <-time.After(2 * time.Second):
        // 触发 OnFailure 事件
        w.InvokeEvent(w.OnFailure, WorkEventArgs{
            Code:         "5000",
            ThreadId:     0,
            WorkName:     w.workName,
            ExceptionMsg: "Timeout adding task",
            SourceEvent:  "Failure",
        })
        return false
    }
}

// Start 启动线程
func (w *Worker) Start() bool {
    var wg sync.WaitGroup
    for i := 0; i < w.threadNumber; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            w.Do()
        }()
    }

    go func() {
        wg.Wait()
        // 触发 OnThreadComplete 事件
        w.InvokeEvent(w.OnThreadComplete, WorkEventArgs{
            Code:     "0",
            ThreadId: 0,
            WorkName: w.workName,
            SourceEvent: "OnThreadComplete",
        })
    }()

    return true
}

// Do 工作方法
func (w *Worker) Do() {
    for {
        select {
        case msg := <-w.workQueue:
            // 触发 OnStart 事件
            w.InvokeEvent(w.OnStart, WorkEventArgs{
                Code:     "0",
                ThreadId: 0,
                WorkName: w.workName,
                SourceEvent: "OnStart",
                Msg:      msg,
            })

            resultChan := make(chan string)
            go func() {
                if w.workerAction != nil {
                    resultChan <- w.workerAction(msg)
                } else {
                    // 如果 workerAction 为 nil，触发 OnFailure 事件
                    w.InvokeEvent(w.OnFailure, WorkEventArgs{
                        Code:         "5001",
                        ThreadId:     0,
                        WorkName:     w.workName,
                        ExceptionMsg: "Worker action is nil",
                        SourceEvent:  "OnFailure",
                        Msg:          msg,
                    })
                    resultChan <- ""
                }
            }()

            select {
            case result := <-resultChan:
                // 触发 OnSuccess 事件
                w.InvokeEvent(w.OnSuccess, WorkEventArgs{
                    Code:     "0",
                    ThreadId: 0,
                    WorkName: w.workName,
                    SourceEvent: "OnSuccess",
                    Msg:      msg,
                    Extend:   result,
                })
            case <-time.After(w.actionTimeout):
                // 触发 OnTimeout 事件
                w.InvokeEvent(w.OnTimeout, WorkEventArgs{
                    Code:         "5003",
                    ThreadId:     0,
                    WorkName:     w.workName,
                    ExceptionMsg: "Timeout",
                    SourceEvent:  "OnTimeout",
                    Msg:          msg,
                })
            }
        case <-w.cts:
            return
        }
        time.Sleep(w.invertalTakeQueue)
    }
}

// InvokeEvent 异步执行事件
func (w *Worker) InvokeEvent(event func(sender interface{}, e WorkEventArgs), e WorkEventArgs) {
    if event == nil {
        return
    }

    if w.asyncInvokeEvent {
        go event(w, e)
    } else {
        event(w, e)
    }
}

// Dispose 销毁对象,正在处理中的线程,将处理完成
func (w *Worker) Dispose() {
    close(w.cts)
    close(w.workQueue)
}

/*
func main() {
    worker := NewWorker(5, "ExampleWorker", func(msg string) string {
        time.Sleep(1 * time.Second)
        return "Processed: " + msg
    }, false, 2*time.Second)

    worker.OnAddedTask = func(sender interface{}, e WorkEventArgs) {
        fmt.Println("Task added:", e.Msg)
    }

    worker.OnSuccess = func(sender interface{}, e WorkEventArgs) {
        fmt.Println("Task succeeded:", e.Msg, e.Extend)
    }

    worker.OnTimeout = func(sender interface{}, e WorkEventArgs) {
        fmt.Println("Task timeout:", e.Msg)
    }

    worker.OnFailure = func(sender interface{}, e WorkEventArgs) {
        fmt.Println("Task failed:", e.Msg, e.ExceptionMsg)
    }

    worker.OnThreadComplete = func(sender interface{}, e WorkEventArgs) {
        fmt.Println("All threads completed")
    }

    worker.Start()

    for i := 0; i < 10; i++ {
        worker.AddTask(fmt.Sprintf("Task %d", i))
    }

    time.Sleep(10 * time.Second)
    worker.Dispose()
}
*/