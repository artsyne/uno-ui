package benchmark

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"log/slog"

	"github.com/biuren/hebe/internal/core/event"
)

// BenchmarkEvent 基准测试事件
type BenchmarkEvent struct {
	Data      string    `json:"data"`
	timestamp time.Time
}

// EventType 返回事件类型
func (e *BenchmarkEvent) EventType() string {
	return "benchmark.event"
}

// Timestamp 返回事件时间戳
func (e *BenchmarkEvent) Timestamp() time.Time {
	return e.timestamp
}

// NewBenchmarkEvent 创建新的基准测试事件
func NewBenchmarkEvent(data string) *BenchmarkEvent {
	return &BenchmarkEvent{
		Data:      data,
		timestamp: time.Now(),
	}
}

// DomainBenchmarkEvent 领域事件基准测试
type DomainBenchmarkEvent struct {
	Data          string    `json:"data"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	Version       int       `json:"version"`
	timestamp     time.Time
}

func (e *DomainBenchmarkEvent) EventType() string {
	return "domain.benchmark"
}

func (e *DomainBenchmarkEvent) Timestamp() time.Time {
	return e.timestamp
}

// SystemBenchmarkEvent 系统事件基准测试
type SystemBenchmarkEvent struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	timestamp time.Time
}

func (e *SystemBenchmarkEvent) EventType() string {
	return "system.benchmark"
}

func (e *SystemBenchmarkEvent) Timestamp() time.Time {
	return e.timestamp
}

// BenchmarkEventBus_Publish 基准测试事件发布
func BenchmarkEventBus_Publish(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	ctx := context.Background()
	testEvent := NewBenchmarkEvent("benchmark data")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bus.Publish(ctx, testEvent)
		}
	})
}

// BenchmarkEventBus_PublishWithHandlers 基准测试带处理器的事件发布
func BenchmarkEventBus_PublishWithHandlers(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	// 添加多个处理器
	for i := 0; i < 10; i++ {
		handler := func(ctx context.Context, e *BenchmarkEvent) error {
			// 模拟一些处理工作
			_ = e.Data
			return nil
		}
		_, _ = bus.Subscribe(handler)
	}

	ctx := context.Background()
	testEvent := NewBenchmarkEvent("benchmark data")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bus.Publish(ctx, testEvent)
		}
	})
}

// BenchmarkEventBus_Subscribe 基准测试事件订阅
func BenchmarkEventBus_Subscribe(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler := func(ctx context.Context, e *BenchmarkEvent) error {
			return nil
		}
		_, _ = bus.Subscribe(handler)
	}
}

// BenchmarkEventBus_Unsubscribe 基准测试事件取消订阅
func BenchmarkEventBus_Unsubscribe(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	// 预先订阅处理器
	subscriptionIDs := make([]event.SubscriptionID, b.N)
	for i := 0; i < b.N; i++ {
		handler := func(ctx context.Context, e *BenchmarkEvent) error {
			return nil
		}
		subscriptionIDs[i], _ = bus.Subscribe(handler)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Unsubscribe(subscriptionIDs[i])
	}
}

// BenchmarkEventBus_ConcurrentPublish 基准测试并发事件发布
func BenchmarkEventBus_ConcurrentPublish(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	// 添加处理器
	var processedCount int64
	handler := func(ctx context.Context, e *BenchmarkEvent) error {
		atomic.AddInt64(&processedCount, 1)
		return nil
	}
	_, _ = bus.Subscribe(handler)

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testEvent := NewBenchmarkEvent(fmt.Sprintf("data-%d", time.Now().UnixNano()))
			bus.Publish(ctx, testEvent)
		}
	})

	// 等待所有事件处理完成
	time.Sleep(100 * time.Millisecond)
}

// BenchmarkEventBus_HighThroughput 基准测试高吞吐量场景
func BenchmarkEventBus_HighThroughput(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	// 添加多个处理器模拟真实场景
	var totalProcessed int64
	for i := 0; i < 5; i++ {
		handler := func(ctx context.Context, e *BenchmarkEvent) error {
			// 模拟一些处理时间
			time.Sleep(time.Microsecond)
			atomic.AddInt64(&totalProcessed, 1)
			return nil
		}
		_, _ = bus.Subscribe(handler)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testEvent := NewBenchmarkEvent(fmt.Sprintf("high-throughput-data-%d", i))
		bus.Publish(ctx, testEvent)
	}

	// 等待所有事件处理完成
	time.Sleep(500 * time.Millisecond)
}

// BenchmarkEventBus_MemoryUsage 基准测试内存使用
func BenchmarkEventBus_MemoryUsage(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	// 添加处理器
	handler := func(ctx context.Context, e *BenchmarkEvent) error {
		return nil
	}
	_, _ = bus.Subscribe(handler)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建大量事件测试内存使用
		testEvent := NewBenchmarkEvent(fmt.Sprintf("memory-test-data-%d-with-some-longer-content-to-test-memory-allocation", i))
		bus.Publish(ctx, testEvent)
	}
}

// BenchmarkTypedEventBusManager 基准测试类型化事件总线管理器
func BenchmarkTypedEventBusManager(b *testing.B) {
	// 直接创建管理器实例，避免全局状态问题
	manager := event.NewTypedEventBusManager[*BenchmarkEvent](slog.Default())
	defer manager.CloseAll()

	// 创建事件总线
	err := manager.CreateBus("benchmark.event", slog.Default())
	if err != nil {
		b.Fatal(err)
	}

	// 获取事件总线
	bus, err := manager.GetBus("benchmark.event")
	if err != nil {
		b.Fatal(err)
	}

	// 添加处理器
	handler := func(ctx context.Context, e *BenchmarkEvent) error {
		return nil
	}
	_, err = bus.Subscribe(handler)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	testEvent := NewBenchmarkEvent("manager data")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bus.Publish(ctx, testEvent)
		}
	})
}

// BenchmarkEventBusManager 基准测试事件总线管理器
func BenchmarkEventBusManager(b *testing.B) {
	// 直接创建管理器实例，避免全局状态问题
	manager := event.NewTypedEventBusManager[*BenchmarkEvent](slog.Default())
	defer manager.CloseAll()

	// 创建多个事件总线
	for i := 0; i < 10; i++ {
		busName := fmt.Sprintf("bus-%d", i)
		err := manager.CreateBus(busName, slog.Default())
		if err != nil {
			b.Fatalf("Failed to create bus %s: %v", busName, err)
		}
	}

	ctx := context.Background()
	testEvent := NewBenchmarkEvent("manager test data")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 随机选择一个事件总线
			busName := fmt.Sprintf("bus-%d", b.N%10)
			bus, err := manager.GetBus(busName)
			if err != nil {
				b.Errorf("Failed to get bus %s: %v", busName, err)
				continue
			}
			bus.Publish(ctx, testEvent)
		}
	})
}

// BenchmarkEventTypes 基准测试不同事件类型
func BenchmarkEventTypes(b *testing.B) {
	ctx := context.Background()

	b.Run("BaseEvent", func(b *testing.B) {
		bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
		defer bus.Close()
		
		handler := func(ctx context.Context, e *BenchmarkEvent) error {
			return nil
		}
		_, _ = bus.Subscribe(handler)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				event := NewBenchmarkEvent("base event data")
				bus.Publish(ctx, event)
			}
		})
	})

	b.Run("DomainEvent", func(b *testing.B) {
		bus := event.NewEventBus[*DomainBenchmarkEvent](slog.Default())
		defer bus.Close()

		handler := func(ctx context.Context, e *DomainBenchmarkEvent) error {
			return nil
		}
		_, _ = bus.Subscribe(handler)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				event := &DomainBenchmarkEvent{
					Data:          "domain event data",
					AggregateID:   "agg1",
					AggregateType: "aggregate",
					Version:       1,
					timestamp:     time.Now(),
				}
				bus.Publish(ctx, event)
			}
		})
	})

	b.Run("SystemEvent", func(b *testing.B) {
		bus := event.NewEventBus[*SystemBenchmarkEvent](slog.Default())
		defer bus.Close()
		ctx := context.Background()

		handler := func(ctx context.Context, e *SystemBenchmarkEvent) error {
			return nil
		}
		_, _ = bus.Subscribe(handler)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				event := &SystemBenchmarkEvent{
					Level:     "info",
					Message:   "system event message",
					Data:      map[string]interface{}{"key": "value"},
					timestamp: time.Now(),
				}
				bus.Publish(ctx, event)
		}
	})
	})
}

// BenchmarkEventBus_ErrorHandling 基准测试错误处理性能
func BenchmarkEventBus_ErrorHandling(b *testing.B) {
	bus := event.NewEventBus[*BenchmarkEvent](slog.Default())
	defer bus.Close()

	// 添加会产生错误的处理器
	errorHandler := func(ctx context.Context, e *BenchmarkEvent) error {
		return fmt.Errorf("handler error")
	}
	_, _ = bus.Subscribe(errorHandler)

	// 添加正常处理器
	normalHandler := func(ctx context.Context, e *BenchmarkEvent) error {
		return nil
	}
	_, _ = bus.Subscribe(normalHandler)

	ctx := context.Background()
	testEvent := NewBenchmarkEvent("error test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 忽略错误，这是错误处理性能测试，不应打印错误日志
		bus.Publish(ctx, testEvent)
	}
}