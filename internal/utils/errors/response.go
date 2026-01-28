package errors

import "time"

// APIErrorResponse 标准化错误响应结构
type APIErrorResponse struct {
	Success   bool        `json:"success" example:"false"`
	Code      int         `json:"code" example:"4001001"`
	Message   string      `json:"message" example:"验证失败"`
	Data      interface{} `json:"data,omitempty"`    // 新增 Data 字段，用于替代 Details
	Details   interface{} `json:"details,omitempty"` // 保留 Details 以向后兼容
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id"`
}

// NewAPIErrorResponse 创建新的标准化错误响应
func NewAPIErrorResponse(code int, message string, details interface{}, requestID string) APIErrorResponse {
	return APIErrorResponse{
		Success:   false,
		Code:      code,
		Message:   message,
		Data:      details, // 同时填充 Data 和 Details
		Details:   details,
		Timestamp: time.Now().UnixMilli(),
		RequestID: requestID,
	}
}
