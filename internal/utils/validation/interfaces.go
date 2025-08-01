package validation

import "context"

// Validator 验证器接口
// 所有需要进行业务逻辑验证的请求对象必须实现此接口
type Validator interface {
	// Validate 执行请求数据的验证
	// 包括结构体标签验证和自定义业务逻辑验证
	// 返回验证错误，如果验证通过返回 nil
	Validate() error
}

// RequestActionGetter 请求操作类型获取接口
// 用于获取请求的操作类型（如 "create", "update", "delete"）
type RequestActionGetter interface {
	// GetAction 获取请求的操作类型
	// 返回操作类型字符串，用于后续的业务逻辑分发
	GetAction() string
}

// EntityPreparer 实体准备器接口（完全去除业务依赖）
// 负责根据操作类型将请求对象转换为对应的领域实体
// 使用泛型 T 确保类型安全
type EntityPreparer[T any] interface {
	// PrepareEntity 根据操作类型准备领域实体
	// 参数：
	//   ctx: 上下文，包含所有必要的环境信息（由调用方通过ContextEnricher注入）
	//   req: 原始请求对象（通常是实现此接口的对象本身）
	//   action: 操作类型（从 RequestActionGetter 获取）
	// 返回：
	//   *T: 准备好的领域实体
	//   error: 准备过程中的错误
	PrepareEntity(ctx context.Context, req interface{}, action string) (*T, error)
}

// RequestProcessor 请求处理器组合接口
// 将所有必要的接口组合在一起，简化使用
type RequestProcessor[T any] interface {
	Validator
	RequestActionGetter
	EntityPreparer[T]
}

// ContextEnricher 上下文增强器接口
// 由调用方实现，用于向context中注入业务相关的信息
// 这样validation包就完全不需要了解具体的业务概念
type ContextEnricher interface {
	// EnrichContext 增强上下文
	// 向context中注入业务相关的信息，如用户ID、时间戳等
	// 返回增强后的context
	EnrichContext(ctx context.Context) context.Context
}

// BusinessHandler 业务处理器函数类型
// 定义业务逻辑处理函数的标准签名
// 接收context和实体，返回处理后的实体和可能的错误
type BusinessHandler[T any] func(ctx context.Context, entity *T) (*T, error)

// BusinessHandlerMap 业务处理器映射类型
// 操作类型到业务处理函数的映射
type BusinessHandlerMap[T any] map[string]BusinessHandler[T]
