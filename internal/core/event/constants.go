// Package event 定义标准事件类型常量
package event

import "regexp"

const (
    // 用户相关事件
    EventTypeUserCreated    = "user.created"
    EventTypeUserUpdated    = "user.updated"
    EventTypeUserDeleted    = "user.deleted"
    EventTypeUserLoggedIn   = "user.logged_in"
    EventTypeUserLoggedOut  = "user.logged_out"
    
    // 任务相关事件
    EventTypeTaskCreated    = "task.created"
    EventTypeTaskCompleted  = "task.completed"
    EventTypeTaskFailed     = "task.failed"
    EmailTaskCreated        = "email.task.created"
    
    // 日记相关事件
    EventTypeDiaryCreated   = "diary.created"
    EventTypeDiaryUpdated   = "diary.updated"
    EventTypeDiaryDeleted   = "diary.deleted"
    
    // 文章相关事件
    EventTypePostCreated    = "post.created"
    EventTypePostUpdated    = "post.updated"
    EventTypePostDeleted    = "post.deleted"
)

var (
    // eventTypePattern 事件类型格式验证正则
    eventTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*\.[a-z][a-z0-9_]*$`)
)

// ValidateEventType 验证事件类型格式
func ValidateEventType(eventType string) bool {
    if eventType == "" {
        return false
    }
    return eventTypePattern.MatchString(eventType)
}