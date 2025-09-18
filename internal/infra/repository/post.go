package repository

import (
	"database/sql"
	"fmt"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/jmoiron/sqlx"
)

// PostRepository 文章数据访问层实现
// 负责处理所有与文章相关的数据库操作
type PostRepository struct {
	db *sqlx.DB
}

// NewPostRepository 创建新的文章数据访问层实例
func NewPostRepository() post.RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}
	return &PostRepository{
		db: db,
	}
}

// Insert 插入文章信息到数据库
func (r *PostRepository) Insert(post *post.Post) (int64, error) {
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

// Update 更新文章信息
func (r *PostRepository) Update(post *post.Post) error {
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

// ListByUserID 根据用户ID和日记本ID分页查询文章信息
func (r *PostRepository) ListByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]post.Post, int, error) {
	infoArray := []post.Post{}
	totalRecords := 0

	// 计算偏移量，注意：页码通常从1开始，所以需要减1
	offset := (pageIndex - 1) * pageSize

	sqlstr := `SELECT * FROM br_post WHERE diary_id = ? AND add_user_id = ? 
		and status = 1 ORDER BY post_id DESC LIMIT ? offset ?`

	err := r.db.Select(&infoArray, sqlstr, diaryID, userID, pageSize, offset)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	//获取总记录数
	totalRecords, err = r.GetCountByUserID(userID, diaryID)

	return infoArray, totalRecords, err
}

// GetCountByUserID 根据用户ID和日记本ID统计信息数量
func (r *PostRepository) GetCountByUserID(userID int64, diaryID int64) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_post WHERE diary_id = ? AND add_user_id = ? and  status = 1 "

	err := r.db.Get(&count, sqlstr, diaryID, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// GetByID 根据文章ID获取文章信息（只返回未删除的文章）
func (r *PostRepository) GetByID(postID int64) (*post.Post, error) {
	sqlstr := "SELECT * FROM br_post WHERE post_id = ? AND status = 1"
	post := &post.Post{}
	err := r.db.Get(post, sqlstr, postID)
	if err == sql.ErrNoRows {
		return nil, errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", postID))
	}
	return post, err
}

// UpdateByID 根据文章ID更新文章信息
func (r *PostRepository) UpdateByID(post *post.Post) error {
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

// DeleteByID 根据文章ID软删除文章信息（设置status=0）
func (r *PostRepository) DeleteByID(postID int64) error {
	sqlstr := "UPDATE br_post SET status = 0 WHERE post_id = ? AND status = 1"
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

// GetPreviousPostByDiaryID 根据当前日记ID和日记本ID获取上一篇日记ID
// 注意：由于列表采用post_id倒序排列，"上一篇"指的是post_id更大的（更新的）日记
func (r *PostRepository) GetPreviousPostByDiaryID(currentPostID int64, diaryID int64) (int64, error) {
	sqlstr := `SELECT post_id FROM br_post 
			   WHERE diary_id = ? AND post_id > ? AND status = 1 
			   ORDER BY post_id ASC LIMIT 1`
	var postID int64
	err := r.db.Get(&postID, sqlstr, diaryID, currentPostID)
	if err == sql.ErrNoRows {
		return 0, nil // 没有找到上一篇日记，返回0而不是错误
	}
	return postID, err
}

// GetNextPostByDiaryID 根据当前日记ID和日记本ID获取下一篇日记ID
// 注意：由于列表采用post_id倒序排列，"下一篇"指的是post_id更小的（更旧的）日记
func (r *PostRepository) GetNextPostByDiaryID(currentPostID int64, diaryID int64) (int64, error) {
	sqlstr := `SELECT post_id FROM br_post 
			   WHERE diary_id = ? AND post_id < ? AND status = 1 
			   ORDER BY post_id DESC LIMIT 1`
	var postID int64
	err := r.db.Get(&postID, sqlstr, diaryID, currentPostID)
	if err == sql.ErrNoRows {
		return 0, nil // 没有找到下一篇日记，返回0而不是错误
	}
	return postID, err
}

// GetStatsByIDs 根据日记ID列表获取统计信息
// 只返回状态为正常(status=1)的日记统计信息
func (r *PostRepository) GetStatsByIDs(postIDs []int64) ([]post.Post, error) {
	if len(postIDs) == 0 {
		return []post.Post{}, nil
	}

	// 构建IN查询的占位符
	query, args, err := sqlx.In(`
		SELECT post_id, diary_id, view_num, like_num, comment_num 
		FROM br_post 
		WHERE post_id IN (?) AND status = 1
		ORDER BY post_id
	`, postIDs)
	if err != nil {
		return nil, err
	}

	// 重新绑定查询
	query = r.db.Rebind(query)

	var posts []post.Post
	err = r.db.Select(&posts, query, args...)
	if err == sql.ErrNoRows {
		return []post.Post{}, nil
	}
	return posts, err
}

// IncrementLikeNum 增加日记的点赞数量
func (r *PostRepository) IncrementLikeNum(postID int64) error {
	sqlstr := "UPDATE br_post SET like_num = like_num + 1 WHERE post_id = ?"
	_, err := r.db.Exec(sqlstr, postID)
	return err
}

// DecrementLikeNum 减少日记的点赞数量
func (r *PostRepository) DecrementLikeNum(postID int64) error {
	sqlstr := "UPDATE br_post SET like_num = like_num - 1 WHERE post_id = ? AND like_num > 0"
	_, err := r.db.Exec(sqlstr, postID)
	return err
}