package repository

import (
	"database/sql"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/jmoiron/sqlx"
)

// CommentRepository 评论数据访问层实现
// 负责处理所有与评论相关的数据库操作
type CommentRepository struct {
	db *sqlx.DB
}

// NewCommentRepository 创建新的评论数据访问层实例
func NewCommentRepository() comment.RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}
	return &CommentRepository{
		db: db,
	}
}

// NewCommentRepositoryWithDB 创建带有指定DB连接的评论仓储实例（用于测试）
func NewCommentRepositoryWithDB(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// Insert 插入评论信息到数据库
func (r *CommentRepository) Insert(comment *comment.Comment) (int64, error) {
	sqlstr := `
        INSERT INTO br_comment(
            entry_id, content, add_utc_time, add_user_id, add_user_ip
        )
        VALUES(?, ?, ?, ?, ?)
    `
	res, err := r.db.Exec(sqlstr, comment.EntryID, comment.Content, comment.AddUtcTime,
		comment.AddUserID, comment.AddUserIP)

	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Delete 删除评论
func (r *CommentRepository) Delete(commentID int64) error {
	sqlstr := `DELETE FROM br_comment WHERE comment_id = ?`
	result, err := r.db.Exec(sqlstr, commentID)
	if err != nil {
		return err
	}

	// 检查是否有行被删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.CommentNotFound.WithDetails(map[string]interface{}{"comment_id": commentID})
	}

	return nil
}

// GetByID 根据评论ID获取评论信息
func (r *CommentRepository) GetByID(commentID int64) (*comment.Comment, error) {
	var commentData comment.Comment
	// JOIN br_user to get handle and display_name
	sqlstr := `SELECT c.comment_id, c.entry_id, c.content, c.add_utc_time, c.add_user_id, c.add_user_ip,
	           u.handle, u.display_name
               FROM br_comment c 
               LEFT JOIN br_user u ON c.add_user_id = u.user_id
               WHERE c.comment_id = ?`
	err := r.db.Get(&commentData, sqlstr, commentID)
	if err == sql.ErrNoRows {
		return nil, errors.CommentNotFound.WithDetails(map[string]interface{}{"comment_id": commentID})
	}
	if err != nil {
		return nil, err
	}
	return &commentData, nil
}

// ListByEntryID 根据日记ID分页查询评论信息
func (r *CommentRepository) ListByEntryID(entryID int64, pageIndex, pageSize int) ([]comment.Comment, int, error) {
	commentArray := []comment.Comment{}
	offset := (pageIndex - 1) * pageSize
	// JOIN br_user to get handle and display_name
	sqlstr := `SELECT c.comment_id, c.entry_id, c.content, c.add_utc_time, c.add_user_id, c.add_user_ip,
	           u.handle, u.display_name
               FROM br_comment c 
               LEFT JOIN br_user u ON c.add_user_id = u.user_id
               WHERE c.entry_id = ?
               ORDER BY c.comment_id DESC LIMIT ? OFFSET ?`
	err := r.db.Select(&commentArray, sqlstr, entryID, pageSize, offset)
	if err == sql.ErrNoRows {
		return []comment.Comment{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}

	count, err := r.GetCountByEntryID(entryID)
	if err != nil {
		return nil, 0, err
	}

	return commentArray, count, nil
}

// GetCountByEntryID 根据日记ID统计评论数量
func (r *CommentRepository) GetCountByEntryID(entryID int64) (int, error) {
	var count int
	sqlstr := `SELECT COUNT(*) FROM br_comment WHERE entry_id = ?`
	err := r.db.Get(&count, sqlstr, entryID)
	return count, err
}

// IncrementEntryCommentNum 增加日记的评论数量
func (r *CommentRepository) IncrementEntryCommentNum(entryID int64) error {
	sqlstr := `UPDATE br_entry SET comment_num = comment_num + 1 WHERE entry_id = ?`
	result, err := r.db.Exec(sqlstr, entryID)
	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}

	return nil
}

// DecrementEntryCommentNum 减少日记的评论数量
func (r *CommentRepository) DecrementEntryCommentNum(entryID int64) error {
	sqlstr := `UPDATE br_entry SET comment_num = GREATEST(comment_num - 1, 0) WHERE entry_id = ?`
	result, err := r.db.Exec(sqlstr, entryID)
	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}

	return nil
}
