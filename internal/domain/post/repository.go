package post

import (
	"database/sql"
	"fmt"

	"github.com/biuren/hebe/internal/core/database"
	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/jmoiron/sqlx"
)

// RepositoryInterface 日记数据访问层接口
// 定义所有与日记相关的数据库操作
type RepositoryInterface interface {
	// Insert 插入日记信息到数据库
	Insert(post *Post) (int64, error)

	// Update 更新日记信息
	Update(post *Post) error

	// ListByUserID 根据用户ID和日记本ID分页查询日记信息
	ListByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Post, int, error)

	// GetCountByUserID 根据用户ID和日记本ID统计信息数量
	GetCountByUserID(userID int64, diaryID int64) (int, error)

	// GetByID 根据日记ID获取日记信息
	GetByID(postID int64) (*Post, error)

	// UpdateByID 根据日记ID更新日记信息
	UpdateByID(post *Post) error

	// DeleteByID 根据日记ID删除日记信息
	DeleteByID(postID int64) error
}

// Repository 日记数据访问层实现
// 负责处理所有与日记相关的数据库操作
type Repository struct {
	db *sqlx.DB
}

// NewRepository 创建新的日记数据访问层实例
func NewRepository() RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}

	return &Repository{
		db: db,
	}
}

// Insert 插入日记信息到数据库
func (r *Repository) Insert(post *Post) (int64, error) {
	sqlstr := `
		INSERT INTO br_post(
			diary_id,is_top,
			title,summary,body,view_num,like_num,comment_num,add_user_id,
			add_user_name,add_user_ip,add_utc_time,update_utc_time,status,info_type
		)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	res, err := r.db.Exec(sqlstr, post.DiaryID,
		post.IsTop, post.Title, post.Summary, post.Body, post.ViewNum, post.LikeNum,
		post.CommentNum, post.AddUserID, post.AddUserName, post.AddUserIP, post.AddUTCTime,
		post.UpdateUTCTime, post.Status, post.InfoType)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// Update 更新日记信息
func (r *Repository) Update(post *Post) error {
	sqlstr := `
		UPDATE br_post SET 
			diary_id = ?, is_top = ?, title = ?, summary = ?, body = ?, 
			view_num = ?, like_num = ?, comment_num = ?, 
			update_utc_time = ?, status = ?, info_type = ? 
		WHERE post_id = ?
	`
	result, err := r.db.Exec(sqlstr, post.DiaryID, post.IsTop, post.Title, post.Summary, post.Body,
		post.ViewNum, post.LikeNum, post.CommentNum, post.UpdateUTCTime, post.Status, post.InfoType,
		post.PostID)
	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", post.PostID))
	}

	return nil
}

// ListByUserID 根据用户ID和日记本ID分页查询日记信息
func (r *Repository) ListByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Post, int, error) {
	infoArray := []Post{}
	totalRecords := 0

	// 计算偏移量，注意：页码通常从1开始，所以需要减1
	offset := (pageIndex - 1) * pageSize

	sqlstr := `SELECT * FROM br_post WHERE diary_id = ? AND add_user_id = ? 
		ORDER BY post_id DESC LIMIT ? offset ?`

	err := r.db.Select(&infoArray, sqlstr, diaryID, userID, pageSize, offset)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	//获取总记录数
	totalRecords, err = r.GetCountByUserID(userID, diaryID)

	return infoArray, totalRecords, err
}

// GetCountByUserID 根据用户ID和日记本ID统计信息数量
func (r *Repository) GetCountByUserID(userID int64, diaryID int64) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_post WHERE diary_id = ? AND add_user_id = ?"

	err := r.db.Get(&count, sqlstr, diaryID, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// GetByID 根据日记ID获取日记信息
func (r *Repository) GetByID(postID int64) (*Post, error) {
	sqlstr := "SELECT * FROM br_post WHERE post_id = ?"
	post := &Post{}
	err := r.db.Get(post, sqlstr, postID)
	if err == sql.ErrNoRows {
		return nil, errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", postID))
	}
	return post, err
}

// UpdateByID 根据日记ID更新日记信息
func (r *Repository) UpdateByID(post *Post) error {
	sqlstr := `
		UPDATE br_post SET 
			diary_id = ?, is_top = ?, title = ?, summary = ?, body = ?, 
			view_num = ?, like_num = ?, comment_num = ?, add_user_id = ?, 
			add_user_name = ?, add_user_ip = ?, add_utc_time = ?, 
			update_utc_time = ?, status = ?, info_type = ? 
		WHERE post_id = ?
	`
	result, err := r.db.Exec(sqlstr, post.DiaryID, post.IsTop, post.Title, post.Summary, post.Body,
		post.ViewNum, post.LikeNum, post.CommentNum, post.AddUserID, post.AddUserName,
		post.AddUserIP, post.AddUTCTime, post.UpdateUTCTime, post.Status, post.InfoType,
		post.PostID)
	if err != nil {
		return err
	}
	
	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", post.PostID))
	}
	
	return nil
}

// DeleteByID 根据日记ID删除日记信息
func (r *Repository) DeleteByID(postID int64) error {
	sqlstr := "DELETE FROM br_post WHERE post_id = ?"
	result, err := r.db.Exec(sqlstr, postID)
	if err != nil {
		return err
	}
	
	// 检查是否有行被删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", postID))
	}
	
	return nil
}