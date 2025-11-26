package repository

import (
	"database/sql"
	"fmt"

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

// Insert 插入评论信息到数据库
func (r *CommentRepository) Insert(comment *comment.Comment) (int64, error) {
	sqlstr := `
        INSERT INTO br_comment(
            post_id, content, add_utc_time, add_user_id, add_user_name, add_user_ip
        )
        VALUES(?, ?, ?, ?, ?, ?)
    `
	res, err := r.db.Exec(sqlstr, comment.PostID, comment.Content, comment.AddUtcTime,
		comment.AddUserID, comment.AddUserName, comment.AddUserIP)

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
		return errors.CommentNotFound.WithDetails(fmt.Sprintf("comment_id: %d", commentID))
	}

	return nil
}

// GetByID 根据评论ID获取评论信息
func (r *CommentRepository) GetByID(commentID int64) (*comment.Comment, error) {
	var commentData comment.Comment
	sqlstr := `SELECT comment_id, post_id, content, add_utc_time, add_user_id, add_user_name, add_user_ip 
               FROM br_comment WHERE comment_id = ?`
	err := r.db.Get(&commentData, sqlstr, commentID)
	if err == sql.ErrNoRows {
		return nil, errors.CommentNotFound.WithDetails(fmt.Sprintf("comment_id: %d", commentID))
	}
	if err != nil {
		return nil, err
	}
	return &commentData, nil
}

// ListByPostID 根据日记ID分页查询评论信息
func (r *CommentRepository) ListByPostID(postID int64, pageIndex, pageSize int) ([]comment.Comment, int, error) {
	commentArray := []comment.Comment{}
	offset := (pageIndex - 1) * pageSize
	sqlstr := `SELECT comment_id, post_id, content, add_utc_time, add_user_id, add_user_name, add_user_ip 
               FROM br_comment WHERE post_id = ?
               ORDER BY comment_id DESC LIMIT ? OFFSET ?`
	err := r.db.Select(&commentArray, sqlstr, postID, pageSize, offset)
	if err == sql.ErrNoRows {
		return []comment.Comment{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}

	count, err := r.GetCountByPostID(postID)
	if err != nil {
		return nil, 0, err
	}

	return commentArray, count, nil
}

// GetCountByPostID 根据日记ID统计评论数量
func (r *CommentRepository) GetCountByPostID(postID int64) (int, error) {
	var count int
	sqlstr := `SELECT COUNT(*) FROM br_comment WHERE post_id = ?`
	err := r.db.Get(&count, sqlstr, postID)
	return count, err
}

// IncrementPostCommentNum 增加日记的评论数量
func (r *CommentRepository) IncrementPostCommentNum(postID int64) error {
	sqlstr := `UPDATE br_post SET comment_num = comment_num + 1 WHERE post_id = ?`
	result, err := r.db.Exec(sqlstr, postID)
	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", postID))
	}

	return nil
}

// DecrementPostCommentNum 减少日记的评论数量
func (r *CommentRepository) DecrementPostCommentNum(postID int64) error {
	sqlstr := `UPDATE br_post SET comment_num = GREATEST(comment_num - 1, 0) WHERE post_id = ?`
	result, err := r.db.Exec(sqlstr, postID)
	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", postID))
	}

	return nil
}
