// Package types 共享类型接口定义
//
// 这个包定义了在多个模块之间共享的核心接口，
// 使用Go 1.24的现代特性，专注于类型安全和扩展性。
package types

import (
	"context"

	"github.com/gin-gonic/gin"
)



// AuthUser 已认证用户信息结构体
type AuthUser struct {
	LoginName    string `form:"login_name" json:"login_name"`
	AccessUrl    string `form:"access_url" json:"access_url"`
	UserId       int64  `form:"user_id" json:"user_id"`
	UserName     string `form:"user_name" json:"user_name"`
	FullName     string `form:"full_name" json:"full_name"`
	Email        string `form:"email" json:"email" db:"email"`
	UserGroup    int8   `form:"user_group" json:"user_group"`
	FaceUrl      string `form:"face_url" json:"face_url"`
	LoginUTCTime int64  `form:"login_utc_time" json:"login_utc_time"`
	TimeZone     string `form:"time_zone" json:"time_zone"`
	Lang         string `form:"lang" json:"lang"`
	Status       int8   `form:"status" json:"status"`
}





// ResponseHandler 响应处理器接口 - 使用泛型确保类型安全
type ResponseHandler[T any] interface {
	JSON(c *gin.Context, data T) error
	JSONSuccess(c *gin.Context, data T) error
	JSONError(c *gin.Context, err error) error
	HTML(c *gin.Context, template string, data T) error
	Redirect(c *gin.Context, url string) error
}



// PageDataBuilder 页面数据构建器接口
type PageDataBuilder[T any] interface {
	Build(ctx context.Context, titleKey string, data T) (*PageData[T], error)
	WithMeta(key, value string) PageDataBuilder[T]
}

// TemplateRenderer 模板渲染器接口
type TemplateRenderer interface {
	Render(ctx context.Context, name string, data any) (string, error)
	RegisterFunction(name string, fn any) error
	Reload() error
}

// ServiceLifecycle 服务生命周期接口
type ServiceLifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	HealthCheck(ctx context.Context) error
}