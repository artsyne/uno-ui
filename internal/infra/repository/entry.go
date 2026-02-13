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

// NewEntryRepositoryWithDB 创建带有指定DB连接的日记仓储实例（用于测试）
func NewEntryRepositoryWithDB(db *sqlx.DB) *EntryRepository {
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
func (r *EntryRepository) ListByUserID(userID int64, diaryID int64, status int8, pageIndex, pageSize int) ([]entry.Entry, int, error) {
	infoArray := []entry.Entry{}
	totalRecords := 0

	// 计算偏移量，注意：页码通常从1开始，所以需要减1
	offset := (pageIndex - 1) * pageSize

	// JOIN br_user to get handle and display_name, JOIN br_diary to get diary_title
	sqlstr := `SELECT e.*, 
		COALESCE(u.handle, '') as handle, 
		COALESCE(u.display_name, '') as display_name,
		d.title as diary_title
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		JOIN br_diary d ON e.diary_id = d.diary_id
		WHERE e.diary_id = ? AND e.add_user_id = ? 
		and e.status = ? ORDER BY e.entry_id DESC LIMIT ? offset ?`

	err := r.db.Select(&infoArray, sqlstr, diaryID, userID, status, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(infoArray) == 0 {
		return []entry.Entry{}, 0, nil
	}

	//获取总记录数
	totalRecords, err = r.GetCountByUserID(userID, diaryID, status)

	return infoArray, totalRecords, err
}

// ListLatestPublicEntriesByUserID 查询用户最新公开日记
func (r *EntryRepository) ListLatestPublicEntriesByUserID(userID int64, limit int) ([]entry.Entry, error) {
	entries := []entry.Entry{}
	sqlstr := `SELECT e.*, 
		COALESCE(u.handle, '') as handle, 
		COALESCE(u.display_name, '') as display_name,
		d.title as diary_title,
		d.add_user_id as diary_owner_id,
		d.social_setting
		FROM br_entry e
		JOIN br_diary d ON e.diary_id = d.diary_id
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		WHERE d.add_user_id = ? 
		AND e.status = 1 
		AND d.status = 1 
		AND d.view_setting = 1 
		ORDER BY e.entry_id DESC 
		LIMIT ?`

	err := r.db.Select(&entries, sqlstr, userID, limit)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// ListByDiaryID 根据日记本ID分页查询日记信息（所有作者）
func (r *EntryRepository) ListByDiaryID(diaryID int64, status int8, pageIndex, pageSize int) ([]entry.Entry, int, error) {
	infoArray := []entry.Entry{}
	totalRecords := 0

	offset := (pageIndex - 1) * pageSize

	// JOIN br_user to get handle and display_name, JOIN br_diary to get diary_title
	sqlstr := `SELECT e.*, 
		COALESCE(u.handle, '') as handle, 
		COALESCE(u.display_name, '') as display_name,
		d.title as diary_title
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		JOIN br_diary d ON e.diary_id = d.diary_id
		WHERE e.diary_id = ? 
		and e.status = ? ORDER BY e.entry_id DESC LIMIT ? offset ?`

	err := r.db.Select(&infoArray, sqlstr, diaryID, status, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(infoArray) == 0 {
		return []entry.Entry{}, 0, nil
	}

	//获取总记录数
	totalRecords, err = r.GetCountByDiaryID(diaryID, status)

	return infoArray, totalRecords, err
}

// GetCountByUserID 根据用户ID和日记本ID统计信息数量
func (r *EntryRepository) GetCountByUserID(userID int64, diaryID int64, status int8) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_entry WHERE diary_id = ? AND add_user_id = ? and  status = ? "

	err := r.db.Get(&count, sqlstr, diaryID, userID, status)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// GetCountByDiaryID 根据日记本ID统计信息数量
func (r *EntryRepository) GetCountByDiaryID(diaryID int64, status int8) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_entry WHERE diary_id = ? and status = ? "

	err := r.db.Get(&count, sqlstr, diaryID, status)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// ListAllByUserID 根据用户ID分页查询日记信息（跨日记本）
func (r *EntryRepository) ListAllByUserID(userID int64, status int8, pageIndex, pageSize int) ([]entry.Entry, int, error) {
	infoArray := []entry.Entry{}
	totalRecords := 0

	offset := (pageIndex - 1) * pageSize

	sqlstr := `SELECT e.*, 
		COALESCE(u.handle, '') as handle, 
		COALESCE(u.display_name, '') as display_name,
		d.title as diary_title
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		JOIN br_diary d ON e.diary_id = d.diary_id
		WHERE e.add_user_id = ? 
		and e.status = ? ORDER BY e.update_utc_time DESC LIMIT ? offset ?`

	err := r.db.Select(&infoArray, sqlstr, userID, status, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(infoArray) == 0 {
		return []entry.Entry{}, 0, nil
	}

	// 获取总记录数
	countSql := "SELECT COUNT(*) FROM br_entry WHERE add_user_id = ? and status = ?"
	err = r.db.Get(&totalRecords, countSql, userID, status)
	if err == sql.ErrNoRows {
		totalRecords = 0
	}

	return infoArray, totalRecords, err
}

// GetByID 根据日记ID获取日记信息（只返回未删除的日记）
func (r *EntryRepository) GetByID(entryID int64) (*entry.Entry, error) {
	// JOIN br_user to get handle and display_name
	sqlstr := `SELECT e.*, u.handle, u.display_name 
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		WHERE e.entry_id = ? AND e.status = ?`
	e := &entry.Entry{}
	err := r.db.Get(e, sqlstr, entryID, entry.StatusActive)
	if err == sql.ErrNoRows {
		return nil, errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}
	return e, err
}

// GetByIDForUser 根据日记ID和用户ID获取日记信息（包含草稿/回收站状态）
func (r *EntryRepository) GetByIDForUser(entryID int64, userID int64) (*entry.Entry, error) {
	sqlstr := `SELECT e.*, COALESCE(u.handle, '') as handle, COALESCE(u.display_name, '') as display_name
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		WHERE e.entry_id = ? AND e.add_user_id = ?`

	e := &entry.Entry{}
	err := r.db.Get(e, sqlstr, entryID, userID)
	if err == sql.ErrNoRows {
		return nil, errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}
	return e, err
}

// GetByIDRaw 根据日记ID获取日记信息（忽略状态和用户限制）
func (r *EntryRepository) GetByIDRaw(entryID int64) (*entry.Entry, error) {
	sqlstr := `SELECT e.*, COALESCE(u.handle, '') as handle, COALESCE(u.display_name, '') as display_name
		FROM br_entry e 
		LEFT JOIN br_user u ON e.add_user_id = u.user_id
		WHERE e.entry_id = ?`

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
	sqlstr := "UPDATE br_entry SET status = ? WHERE entry_id = ?"
	result, err := r.db.Exec(sqlstr, entry.StatusDeleted, entryID)
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

// HardDeleteByID 根据日记ID物理删除日记信息
func (r *EntryRepository) HardDeleteByID(entryID int64) error {
	sqlstr := "DELETE FROM br_entry WHERE entry_id = ?"
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
			   WHERE diary_id = ? AND entry_id > ? AND status = ? 
			   ORDER BY entry_id ASC LIMIT 1`
	var entryID int64
	err := r.db.Get(&entryID, sqlstr, diaryID, currentEntryID, entry.StatusActive)
	if err == sql.ErrNoRows {
		return 0, nil // 没有找到上一篇日记，返回0而不是错误
	}
	return entryID, err
}

// GetNextEntryByDiaryID 根据当前日记ID和日记本ID获取下一篇日记ID
// 注意：由于列表采用entry_id倒序排列，"下一篇"指的是entry_id更小的（更旧的）日记
func (r *EntryRepository) GetNextEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
	sqlstr := `SELECT entry_id FROM br_entry 
			   WHERE diary_id = ? AND entry_id < ? AND status = ? 
			   ORDER BY entry_id DESC LIMIT 1`
	var entryID int64
	err := r.db.Get(&entryID, sqlstr, diaryID, currentEntryID, entry.StatusActive)
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
		WHERE entry_id IN (?) AND status = ?
		ORDER BY entry_id
	`, entryIDs, entry.StatusActive)
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
