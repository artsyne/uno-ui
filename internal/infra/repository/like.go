package repository

import (
	"fmt"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/jmoiron/sqlx"
)

// LikeRepository 点赞数据访问层实现
type LikeRepository struct {
	db *sqlx.DB
}

// NewLikeRepository 创建点赞仓库实例
func NewLikeRepository() like.Repository {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}
	return &LikeRepository{
		db: db,
	}
}

// NewLikeRepositoryWithDB 创建带有指定DB连接的点赞仓储实例（用于测试）
func NewLikeRepositoryWithDB(db *sqlx.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

// ToggleLike 切换点赞状态
func (r *LikeRepository) ToggleLike(userID, ownerID int64, likeType int) error {
	// 检查是否已存在点赞记录
	var count int
	checkSQL := `SELECT COUNT(*) FROM br_like WHERE user_id = ? AND owner_id = ? AND like_type = ?`
	err := r.db.QueryRow(checkSQL, userID, ownerID, likeType).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查点赞状态失败: %w", err)
	}

	if count > 0 {
		// 已存在，删除点赞记录
		deleteSQL := `DELETE FROM br_like WHERE user_id = ? AND owner_id = ? AND like_type = ?`
		_, err = r.db.Exec(deleteSQL, userID, ownerID, likeType)
		if err != nil {
			return fmt.Errorf("取消点赞失败: %w", err)
		}
	} else {
		// 不存在，插入点赞记录
		insertSQL := `INSERT INTO br_like (user_id, owner_id, like_type) VALUES (?, ?, ?)`
		_, err = r.db.Exec(insertSQL, userID, ownerID, likeType)
		if err != nil {
			return fmt.Errorf("添加点赞失败: %w", err)
		}
	}

	return nil
}

// IsLiked 检查用户是否已点赞
func (r *LikeRepository) IsLiked(userID, ownerID int64, likeType int) (bool, error) {
	var count int
	checkSQL := `SELECT COUNT(*) FROM br_like WHERE user_id = ? AND owner_id = ? AND like_type = ?`
	err := r.db.QueryRow(checkSQL, userID, ownerID, likeType).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	return count > 0, nil
}

// GetLikeCount 获取点赞总数（从 br_like 表查询）
// 注意：推荐通过 owner 对象的 likeNum 字段获取，而不是使用此方法
func (r *LikeRepository) GetLikeCount(ownerID int64, likeType int) (int, error) {
	var count int
	countSQL := `SELECT COUNT(*) FROM br_like WHERE owner_id = ? AND like_type = ?`
	err := r.db.QueryRow(countSQL, ownerID, likeType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("获取点赞数量失败: %w", err)
	}

	return count, nil
}