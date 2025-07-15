package validation

import (
	"context"

	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// HandleRequest 通用请求处理器（完全业务无关）
// 这是框架的核心函数，负责统一处理所有类型的请求
//
// 处理流程：
// 1. 验证请求数据（调用 processor.Validate()）
// 2. 获取操作类型（调用 processor.GetAction()）
// 3. 增强上下文（调用 contextEnricher.EnrichContext()）
// 4. 准备领域实体（调用 processor.PrepareEntity()）
// 5. 查找并执行业务处理函数
// 6. 返回处理结果
//
// 类型参数：
//
//	T: 领域实体类型（如 Diary, Post 等）
//
// 参数：
//
//	ginCtx: Gin 上下文，用于日志记录和错误处理
//	processor: 请求处理器，实现了所有必要的接口
//	businessHandlers: 业务处理函数映射
//	contextEnricher: 上下文增强器（可选，用于注入业务信息）
//
// 返回值：
//
//	*T: 处理后的领域实体
//	error: 处理过程中的错误
func HandleRequest[T any](
	ginCtx *gin.Context,
	processor RequestProcessor[T],
	businessHandlers BusinessHandlerMap[T],
	contextEnricher ContextEnricher,
) (*T, error) {
	// 第一步：验证请求数据
	// 包括结构体标签验证和自定义业务逻辑验证
	if err := processor.Validate(); err != nil {
		return nil, err
	}

	// 第二步：获取操作类型
	// 用于后续的业务逻辑分发
	action := processor.GetAction()

	// 第三步：创建并增强上下文
	// 从gin context创建基础context
	ctx := context.Background()

	// 如果提供了上下文增强器，则增强上下文
	if contextEnricher != nil {
		ctx = contextEnricher.EnrichContext(ctx)
	}

	// 第四步：准备领域实体
	// 根据操作类型将请求对象转换为领域实体
	// 所有业务相关的信息都从context中获取
	entity, err := processor.PrepareEntity(ctx, processor, action)
	if err != nil {
		return nil, err
	}

	// 第五步：查找业务处理函数
	// 根据操作类型从映射中获取对应的处理函数
	handlerFunc, exists := businessHandlers[action]
	if !exists {
		// 如果找不到对应的处理函数，返回错误
		// 通常这种情况不会发生，因为在验证阶段已经检查了 action 的有效性
		return nil, errors.InvalidParams.WithDetails("不支持的操作类型：" + action)
	}

	// 第六步：执行业务逻辑
	// 调用具体的业务处理函数，传入context和实体
	return handlerFunc(ctx, entity)
}
