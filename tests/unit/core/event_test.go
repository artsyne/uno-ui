package core

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/core/event"
)

// TestEvent 测试用的事件类型
type TestEvent struct {
	Data      string    `json:"data"`
	timestamp time.Time
}

// EventType 实现 Event 接口
func (e *TestEvent) EventType() string {
	return "test.event"
}

// Timestamp 实现 Event 接口
func (e *TestEvent) Timestamp() time.Time {
	return e.timestamp
}

// NewTestEvent 创建测试事件
func NewTestEvent(data string) *TestEvent {
	return &TestEvent{
		Data:      data,
		timestamp: time.Now(),
	}
}

// AnotherTestEvent 另一个测试事件类型
type AnotherTestEvent struct {
	Value     int       `json:"value"`
	timestamp time.Time
}

// EventType 实现 Event 接口
func (e *AnotherTestEvent) EventType() string {
	return "another.test.event"
}

// Timestamp 实现 Event 接口
func (e *AnotherTestEvent) Timestamp() time.Time {
	return e.timestamp
}

// NewAnotherTestEvent 创建另一个测试事件
func NewAnotherTestEvent(value int) *AnotherTestEvent {
	return &AnotherTestEvent{
		Value:     value,
		timestamp: time.Now(),
	}
}



// TestEventBus_NewEventBus 测试事件总线创建
func TestEventBus_NewEventBus(t *testing.T) {
	// 测试泛型事件总线创建
	bus := event.NewEventBus[*TestEvent](slog.Default())
	if bus == nil {
		t.Fatal("Expected non-nil event bus")
	}

	// 测试类型化事件总线管理器创建
	manager := event.GetTypedEventBusManager[*TestEvent]("test", slog.Default())
	if manager == nil {
		t.Fatal("Expected non-nil typed event bus manager")
	}
	defer manager.CloseAll()
}

// TestEventBus_Subscribe 测试事件订阅
func TestEventBus_Subscribe(t *testing.T) {
	bus := event.NewEventBus[*TestEvent](slog.Default())
	defer bus.Close()

	// 订阅处理器
	handler := func(ctx context.Context, e *TestEvent) error {
		return nil
	}

	subscriptionID, err := bus.Subscribe(handler)
	if err != nil {
		t.Fatalf("Failed to subscribe handler: %v", err)
	}

	if subscriptionID == 0 {
		t.Error("Expected non-zero subscription ID")
	}
}

// TestEventBus_Publish 测试事件发布
func TestEventBus_Publish(t *testing.T) {
	bus := event.NewEventBus[*TestEvent](slog.Default())
	defer bus.Close()
	ctx := context.Background()

	// 记录处理的事件
	var handledEvents []*TestEvent
	var mu sync.Mutex

	// 订阅处理器
	handler := func(ctx context.Context, e *TestEvent) error {
		mu.Lock()
		defer mu.Unlock()
		handledEvents = append(handledEvents, e)
		return nil
	}

	_, err := bus.Subscribe(handler)
	if err != nil {
		t.Fatalf("Failed to subscribe handler: %v", err)
	}

	// 发布事件
	testEvent := NewTestEvent("test data")
	err = bus.Publish(ctx, testEvent)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证事件被处理
	mu.Lock()
	if len(handledEvents) != 1 {
		t.Fatalf("Expected 1 handled event, got %d", len(handledEvents))
	}
	if handledEvents[0].Data != "test data" {
		t.Errorf("Expected event data 'test data', got '%s'", handledEvents[0].Data)
	}
	mu.Unlock()
}

// TestEventBus_MultipleHandlers 测试多个处理器
func TestEventBus_MultipleHandlers(t *testing.T) {
	bus := event.NewEventBus[*TestEvent](slog.Default())
	defer bus.Close()
	ctx := context.Background()

	// 计数器
	var count1, count2 int32

	// 第一个处理器
	handler1 := func(ctx context.Context, e *TestEvent) error {
		atomic.AddInt32(&count1, 1)
		return nil
	}

	// 第二个处理器
	handler2 := func(ctx context.Context, e *TestEvent) error {
		atomic.AddInt32(&count2, 1)
		return nil
	}

	// 订阅处理器
	_, err := bus.Subscribe(handler1)
	if err != nil {
		t.Fatalf("Failed to subscribe handler1: %v", err)
	}

	_, err = bus.Subscribe(handler2)
	if err != nil {
		t.Fatalf("Failed to subscribe handler2: %v", err)
	}

	// 发布事件
	testEvent := NewTestEvent("test data")
	err = bus.Publish(ctx, testEvent)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证两个处理器都被调用
	if atomic.LoadInt32(&count1) != 1 {
		t.Errorf("Expected handler1 to be called 1 time, got %d", count1)
	}
	if atomic.LoadInt32(&count2) != 1 {
		t.Errorf("Expected handler2 to be called 1 time, got %d", count2)
	}
}

// TestEventBus_Unsubscribe 测试取消订阅
func TestEventBus_Unsubscribe(t *testing.T) {
	bus := event.NewEventBus[*TestEvent](slog.Default())
	defer bus.Close()
	ctx := context.Background()

	// 计数器
	var count int32

	// 处理器
	handler := func(ctx context.Context, e *TestEvent) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	// 订阅处理器
	subscriptionID, err := bus.Subscribe(handler)
	if err != nil {
		t.Fatalf("Failed to subscribe handler: %v", err)
	}

	// 发布第一个事件
	testEvent1 := NewTestEvent("test data 1")
	err = bus.Publish(ctx, testEvent1)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 取消订阅
	err = bus.Unsubscribe(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to unsubscribe handler: %v", err)
	}

	// 发布第二个事件
	testEvent2 := NewTestEvent("test data 2")
	err = bus.Publish(ctx, testEvent2)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证只处理了第一个事件
	if atomic.LoadInt32(&count) != 1 {
		t.Errorf("Expected handler to be called 1 time, got %d", count)
	}

	// 测试取消不存在的订阅
	err = bus.Unsubscribe(999)
	if err == nil {
		t.Error("Expected error when unsubscribing non-existent subscription")
	}
}

// TestEventBus_ErrorHandling 测试错误处理
func TestEventBus_ErrorHandling(t *testing.T) {
	bus := event.NewEventBus[*TestEvent](slog.Default())
	defer bus.Close()
	ctx := context.Background()

	// 错误处理器
	errorHandler := func(ctx context.Context, e *TestEvent) error {
		return errors.New("handler error")
	}

	// 正常处理器
	var normalHandlerCalled bool
	normalHandler := func(ctx context.Context, e *TestEvent) error {
		normalHandlerCalled = true
		return nil
	}

	// 订阅处理器
	_, err := bus.Subscribe(errorHandler)
	if err != nil {
		t.Fatalf("Failed to subscribe error handler: %v", err)
	}

	_, err = bus.Subscribe(normalHandler)
	if err != nil {
		t.Fatalf("Failed to subscribe normal handler: %v", err)
	}

	// 发布事件
	testEvent := NewTestEvent("test data")
	err = bus.Publish(ctx, testEvent)
	if err == nil {
		t.Error("Expected error from publish, got nil")
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证正常处理器仍然被调用
	if !normalHandlerCalled {
		t.Error("Expected normal handler to be called")
	}


}

// TestEventBus_PanicRecovery 测试 panic 恢复
func TestEventBus_PanicRecovery(t *testing.T) {
	bus := event.NewEventBus[*TestEvent](slog.Default())
	defer bus.Close()
	ctx := context.Background()

	// Panic 处理器
	panicHandler := func(ctx context.Context, e *TestEvent) error {
		panic("handler panic")
	}

	// 正常处理器
	var normalHandlerCalled bool
	normalHandler := func(ctx context.Context, e *TestEvent) error {
		normalHandlerCalled = true
		return nil
	}

	// 订阅处理器
	_, err := bus.Subscribe(panicHandler)
	if err != nil {
		t.Fatalf("Failed to subscribe panic handler: %v", err)
	}

	_, err = bus.Subscribe(normalHandler)
	if err != nil {
		t.Fatalf("Failed to subscribe normal handler: %v", err)
	}

	// 发布事件
	testEvent := NewTestEvent("test data")
	err = bus.Publish(ctx, testEvent)
	if err == nil {
		t.Error("Expected error from publish due to panic, got nil")
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证正常处理器仍然被调用
	if !normalHandlerCalled {
		t.Error("Expected normal handler to be called")
	}


}

// TestEventBus_Close 测试关闭事件总线
func TestEventBus_Close(t *testing.T) {
	bus := event.NewEventBus[*TestEvent](slog.Default())
	ctx := context.Background()

	// 关闭事件总线
	err := bus.Close()
	if err != nil {
		t.Fatalf("Failed to close event bus: %v", err)
	}

	// 尝试订阅（应该失败）
	handler := func(ctx context.Context, e *TestEvent) error {
		return nil
	}
	_, err = bus.Subscribe(handler)
	if err == nil {
		t.Error("Expected error when subscribing to closed bus")
	}

	// 尝试发布事件（应该失败）
	testEvent := NewTestEvent("test data")
	err = bus.Publish(ctx, testEvent)
	if err == nil {
		t.Error("Expected error when publishing to closed bus")
	}

	// 尝试再次关闭（应该失败）
	err = bus.Close()
	if err == nil {
		t.Error("Expected error when closing already closed bus")
	}
}

// TestEventBus_DifferentEventTypes 测试不同事件类型
func TestEventBus_DifferentEventTypes(t *testing.T) {
	ctx := context.Background()

	// 计数器
	var testEventCount, anotherEventCount int32

	// TestEvent 处理器
	testHandler := func(ctx context.Context, e *TestEvent) error {
		atomic.AddInt32(&testEventCount, 1)
		return nil
	}

	// AnotherTestEvent 处理器
	anotherHandler := func(ctx context.Context, e *AnotherTestEvent) error {
		atomic.AddInt32(&anotherEventCount, 1)
		return nil
	}

	// 创建不同类型的事件总线并订阅处理器
	testBus := event.NewEventBus[*TestEvent](slog.Default())
	defer testBus.Close()
	_, err := testBus.Subscribe(testHandler)
	if err != nil {
		t.Fatalf("Failed to subscribe test handler: %v", err)
	}

	anotherBus := event.NewEventBus[*AnotherTestEvent](slog.Default())
	defer anotherBus.Close()
	_, err = anotherBus.Subscribe(anotherHandler)
	if err != nil {
		t.Fatalf("Failed to subscribe another handler: %v", err)
	}

	// 发布不同类型的事件
	testEvent := NewTestEvent("test data")
	err = testBus.Publish(ctx, testEvent)
	if err != nil {
		t.Fatalf("Failed to publish test event: %v", err)
	}

	anotherEvent := NewAnotherTestEvent(42)
	err = anotherBus.Publish(ctx, anotherEvent)
	if err != nil {
		t.Fatalf("Failed to publish another event: %v", err)
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证各自的处理器被调用
	if atomic.LoadInt32(&testEventCount) != 1 {
		t.Errorf("Expected test event handler to be called 1 time, got %d", testEventCount)
	}
	if atomic.LoadInt32(&anotherEventCount) != 1 {
		t.Errorf("Expected another event handler to be called 1 time, got %d", anotherEventCount)
	}
}