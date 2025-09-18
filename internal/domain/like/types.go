package like

// Like 点赞结构体
type Like struct {
	LikeID   int64 `json:"like_id" db:"like_id"`     // 点赞ID
	UserID   int64 `json:"user_id" db:"user_id"`     // 用户ID
	OwnerID  int64 `json:"owner_id" db:"owner_id"`   // 被点赞对象ID（日记本ID或日记ID）
	LikeType int   `json:"like_type" db:"like_type"` // 点赞类型：1=日记本点赞，2=日记点赞
}

// 点赞类型常量
const (
	LikeTypeDiary = 1 // 日记本点赞
	LikeTypePost  = 2 // 日记点赞
)