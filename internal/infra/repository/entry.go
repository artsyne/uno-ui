package repository

import (
	"database/sql"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/jmoiron/sqlx"
)

// EntryRepository 日记数据访问层实现
// 负责处理所有与日记相关的数据库操作
type EntryRepository struct {
	db *sqlx.DB
}

// NewEntryRepository 创建新的日记数据访问层实例
func NewEntryRepository() entry.RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}
	return &EntryRepository{
		db: db,
	}
}

// Insert 插入日记信息到数据库
func (r *EntryRepository) Insert(e *entry.Entry) (int64, error) {
	sqlstr := `
		INSERT INTO br_entry(
			diary_id,is_top,
			title,summary,body,view_num,like_num,comment_num,add_user_id,
			add_user_ip,add_utc_time,update_utc_time,status,info_type
		)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	res, err := r.db.Exec(sqlstr, e.DiaryID,
		e.IsTop, e.Title, e.Summary, e.Body, e.ViewNum, e.LikeNum,
		e.CommentNum, e.AddUserID, e.AddUserIP, e.AddUTCTime,
		e.UpdateUTCTime, e.Status, e.InfoType)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// Update 更新日记信息
func (r *EntryRepository) Update(e *entry.Entry) error {
	sqlstr := `
		UPDATE br_entry SET 
			diary_id = ?, is_top = ?, title = ?, summary = ?, body = ?, 
			view_num = ?, like_num = ?, comment_num = ?, 
			update_utc_time = ?, status = ?, info_type = ? 
		WHERE entry_id = ?
	`
	result, err := r.db.Exec(sqlstr, e.DiaryID, e.IsTop, e.Title, e.Summary, e.Body,
		e.ViewNum, e.LikeNum, e.CommentNum, e.UpdateUTCTime, e.Status, e.InfoType,
		e.EntryID)
	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": e.EntryID})
	}

	return nil
}

// ListByUserID 根据用户ID和日记本ID分页查询日记信息
func (r *EntryRepository) ListByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]entry.Entry, int, error) {
	infoArray := []entry.Entry{}
	totalRecords := 0

	// 计算偏移量，注意：页码通常从1开始，所以需要减1
	offset := (pageIndex - 1) * pageSize

	// JOIN br_user to get handle and display_name
	sqlstr := `SELECT e.*, u.handle, u.display_name 
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		WHERE e.diary_id = ? AND e.add_user_id = ? 
		and e.status = 1 ORDER BY e.entry_id DESC LIMIT ? offset ?`

	err := r.db.Select(&infoArray, sqlstr, diaryID, userID, pageSize, offset)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	//获取总记录数
	totalRecords, err = r.GetCountByUserID(userID, diaryID)

	return infoArray, totalRecords, err
}

// GetCountByUserID 根据用户ID和日记本ID统计信息数量
func (r *EntryRepository) GetCountByUserID(userID int64, diaryID int64) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_entry WHERE diary_id = ? AND add_user_id = ? and  status = 1 "

	err := r.db.Get(&count, sqlstr, diaryID, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// GetByID 根据日记ID获取日记信息（只返回未删除的日记）
func (r *EntryRepository) GetByID(entryID int64) (*entry.Entry, error) {
	// JOIN br_user to get handle and display_name
	sqlstr := `SELECT e.*, u.handle, u.display_name 
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		WHERE e.entry_id = ? AND e.status = 1`
	e := &entry.Entry{}
	err := r.db.Get(e, sqlstr, entryID)
	if err == sql.ErrNoRows {
		return nil, errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}
	return e, err
}

// UpdateByID 根据日记ID更新日记信息
func (r *EntryRepository) UpdateByID(e *entry.Entry) error {
	sqlstr := `
		UPDATE br_entry SET 
			diary_id = ?, is_top = ?, title = ?, summary = ?, body = ?, 
			view_num = ?, like_num = ?, comment_num = ?, add_user_id = ?, 
			add_user_ip = ?, add_utc_time = ?, 
			update_utc_time = ?, status = ?, info_type = ? 
		WHERE entry_id = ?
	`
	result, err := r.db.Exec(sqlstr, e.DiaryID, e.IsTop, e.Title, e.Summary, e.Body,
		e.ViewNum, e.LikeNum, e.CommentNum, e.AddUserID,
		e.AddUserIP, e.AddUTCTime, e.UpdateUTCTime, e.Status, e.InfoType,
		e.EntryID)
	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": e.EntryID})
	}

	return nil
}

// DeleteByID 根据日记ID软删除日记信息（设置status=0）
func (r *EntryRepository) DeleteByID(entryID int64) error {
	sqlstr := "UPDATE br_entry SET status = 0 WHERE entry_id = ? AND status = 1"
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

// GetPreviousEntryByDiaryID 根据当前日记ID和日记本ID获取上一篇日记ID
// 注意：由于列表采用entry_id倒序排列，"上一篇"指的是entry_id更大的（更新的）日记
func (r *EntryRepository) GetPreviousEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
	sqlstr := `SELECT entry_id FROM br_entry 
			   WHERE diary_id = ? AND entry_id > ? AND status = 1 
			   ORDER BY entry_id ASC LIMIT 1`
	var entryID int64
	err := r.db.Get(&entryID, sqlstr, diaryID, currentEntryID)
	if err == sql.ErrNoRows {
		return 0, nil // 没有找到上一篇日记，返回0而不是错误
	}
	return entryID, err
}

// GetNextEntryByDiaryID 根据当前日记ID和日记本ID获取下一篇日记ID
// 注意：由于列表采用entry_id倒序排列，"下一篇"指的是entry_id更小的（更旧的）日记
func (r *EntryRepository) GetNextEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
	sqlstr := `SELECT entry_id FROM br_entry 
			   WHERE diary_id = ? AND entry_id < ? AND status = 1 
			   ORDER BY entry_id DESC LIMIT 1`
	var entryID int64
	err := r.db.Get(&entryID, sqlstr, diaryID, currentEntryID)
	if err == sql.ErrNoRows {
		return 0, nil // 没有找到下一篇日记，返回0而不是错误
	}
	return entryID, err
}

// GetStatsByIDs 根据日记ID列表获取统计信息
// 只返回状态为正常(status=1)的日记统计信息
func (r *EntryRepository) GetStatsByIDs(entryIDs []int64) ([]entry.Entry, error) {
	if len(entryIDs) == 0 {
		return []entry.Entry{}, nil
	}

	// 构建IN查询的占位符
	query, args, err := sqlx.In(`
		SELECT entry_id, diary_id, view_num, like_num, comment_num 
		FROM br_entry 
		WHERE entry_id IN (?) AND status = 1
		ORDER BY entry_id
	`, entryIDs)
	if err != nil {
		return nil, err
	}

	// 重新绑定查询
	query = r.db.Rebind(query)

	var entries []entry.Entry
	err = r.db.Select(&entries, query, args...)
	if err == sql.ErrNoRows {
		return []entry.Entry{}, nil
	}
	return entries, err
}

// IncrementLikeNum 增加日记的点赞数量
func (r *EntryRepository) IncrementLikeNum(entryID int64) error {
	sqlstr := "UPDATE br_entry SET like_num = like_num + 1 WHERE entry_id = ?"
	_, err := r.db.Exec(sqlstr, entryID)
	return err
}

// DecrementLikeNum 减少日记的点赞数量
func (r *EntryRepository) DecrementLikeNum(entryID int64) error {
	sqlstr := "UPDATE br_entry SET like_num = like_num - 1 WHERE entry_id = ? AND like_num > 0"
	_, err := r.db.Exec(sqlstr, entryID)
	return err
}
