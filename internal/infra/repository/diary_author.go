package repository

import (
	"database/sql"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/utils/errors"
)

// InsertAuthor 新增日记本作者
func (r *DiaryRepository) InsertAuthor(author *diary.Author) error {
	sqlstr := `
		INSERT INTO br_diary_author
		(diary_id, user_id, is_admin, join_utc_time, post_num, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(sqlstr, author.DiaryID, author.UserID, author.IsAdmin, author.JoinUTCTime, author.PostNum, author.Status)
	return err
}

// UpdateAuthorStatus 更新作者状态
func (r *DiaryRepository) UpdateAuthorStatus(diaryID, userID int64, status int8) error {
	sqlstr := "UPDATE br_diary_author SET status = ? WHERE diary_id = ? AND user_id = ?"
	_, err := r.db.Exec(sqlstr, status, diaryID, userID)
	return err
}

// UpdateAuthorAdmin 更新作者管理员权限
func (r *DiaryRepository) UpdateAuthorAdmin(diaryID, userID int64, isAdmin int8) error {
	sqlstr := "UPDATE br_diary_author SET is_admin = ? WHERE diary_id = ? AND user_id = ?"
	_, err := r.db.Exec(sqlstr, isAdmin, diaryID, userID)
	return err
}

// DeleteAuthor 删除作者
func (r *DiaryRepository) DeleteAuthor(diaryID, userID int64) error {
	sqlstr := "DELETE FROM br_diary_author WHERE diary_id = ? AND user_id = ?"
	_, err := r.db.Exec(sqlstr, diaryID, userID)
	return err
}

// ListAuthorsByDiaryID 获取日记本的所有作者
// onlyActive: true=仅获取活跃作者, false=获取所有状态
// pageIndex, pageSize: 分页参数，如果 pageSize <= 0 则返回所有
func (r *DiaryRepository) ListAuthorsByDiaryID(diaryID int64, onlyActive bool, pageIndex, pageSize int) ([]diary.Author, int, error) {
	var authors []diary.Author
	var err error
	var args []interface{}

	// Base query
	baseQuery := `
			FROM br_diary_author da
			JOIN br_user u ON da.user_id = u.user_id
			WHERE da.diary_id = ?
		`
	args = append(args, diaryID)

	if onlyActive {
		baseQuery += ` AND da.status = ?`
		args = append(args, diary.AuthorStatusActive)
	}

	// Count query
	var count int
	countQuery := `SELECT count(*) ` + baseQuery
	err = r.db.Get(&count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `SELECT da.*, u.handle, u.display_name, u.face_url ` + baseQuery

	// Pagination
	if pageSize > 0 {
		if pageIndex < 1 {
			pageIndex = 1
		}
		offset := (pageIndex - 1) * pageSize
		dataQuery += ` LIMIT ? OFFSET ?`
		args = append(args, pageSize, offset)
	}

	err = r.db.Select(&authors, dataQuery, args...)
	if err == sql.ErrNoRows {
		return []diary.Author{}, 0, nil
	}
	return authors, count, err
}

// ListJoinedDiariesByUserID 获取用户参与的日记本列表（状态为Active）
func (r *DiaryRepository) ListJoinedDiariesByUserID(userID int64) ([]diary.Diary, error) {
	var diaries []diary.Diary
	sqlstr := `
		SELECT d.*, u.handle, u.display_name
		FROM br_diary d
		JOIN br_diary_author da ON d.diary_id = da.diary_id
		LEFT JOIN br_user u ON d.add_user_id = u.user_id
		WHERE da.user_id = ? AND da.status = ? AND d.status = ? AND d.add_user_id != ?
		ORDER BY d.diary_id DESC
	`
	err := r.db.Select(&diaries, sqlstr, userID, diary.AuthorStatusActive, diary.StatusActive, userID)
	if err == sql.ErrNoRows {
		return []diary.Diary{}, nil
	}
	return diaries, err
}

// ListJoinedDiariesByUserIDWithLimit 获取用户参与的日记本列表（带限制，不查询总数）
// 包含 Pending (Status=2) 并按 Status DESC 排序
func (r *DiaryRepository) ListJoinedDiariesByUserIDWithLimit(userID int64, includePrivate bool, limit int) ([]diary.Diary, error) {
	var diaries []diary.Diary
	baseSQL := `
		SELECT d.*, u.handle, u.display_name, da.status as author_status
		FROM br_diary d
		JOIN br_diary_author da ON d.diary_id = da.diary_id
		LEFT JOIN br_user u ON d.add_user_id = u.user_id
		WHERE da.user_id = ?`

	var args []interface{}
	args = append(args, userID)

	if includePrivate {
		// Owner can see Active and Pending
		baseSQL += " AND da.status IN (?, ?)"
		args = append(args, diary.AuthorStatusActive, diary.AuthorStatusPending)
	} else {
		// Others only see Active
		baseSQL += " AND da.status = ?"
		args = append(args, diary.AuthorStatusActive)
	}

	baseSQL += " AND d.status = ? AND d.add_user_id != ?"
	args = append(args, diary.StatusActive, userID)

	if !includePrivate {
		baseSQL += " AND d.view_setting = ?"
		args = append(args, diary.ViewSettingPublic)
	}

	baseSQL += ` ORDER BY da.status DESC, da.join_utc_time DESC, d.diary_id DESC LIMIT ?`
	args = append(args, limit)

	err := r.db.Select(&diaries, baseSQL, args...)
	if err == sql.ErrNoRows {
		return []diary.Diary{}, nil
	}
	if err != nil {
		return nil, err
	}
	return diaries, nil
}

// ListJoinedDiariesByUserIDPaged 获取用户参与的分页日记本列表
// Modified to include Pending (Status=2) and sort by Status DESC (Pending first)
func (r *DiaryRepository) ListJoinedDiariesByUserIDPaged(userID int64, includePrivate bool, pageIndex, pageSize int) ([]diary.Diary, int, error) {
	var diaries []diary.Diary
	offset := (pageIndex - 1) * pageSize

	baseSQL := `
		SELECT d.*, u.handle, u.display_name, da.status as author_status
		FROM br_diary d
		JOIN br_diary_author da ON d.diary_id = da.diary_id
		LEFT JOIN br_user u ON d.add_user_id = u.user_id
		WHERE da.user_id = ?`

	var args []interface{}
	args = append(args, userID)

	if includePrivate {
		baseSQL += " AND da.status IN (?, ?)"
		args = append(args, diary.AuthorStatusActive, diary.AuthorStatusPending)
	} else {
		baseSQL += " AND da.status = ?"
		args = append(args, diary.AuthorStatusActive)
	}

	baseSQL += " AND d.status = ? AND d.add_user_id != ?"
	args = append(args, diary.StatusActive, userID)

	if !includePrivate {
		baseSQL += " AND d.view_setting = ?"
		args = append(args, diary.ViewSettingPublic)
	}

	baseSQL += ` ORDER BY da.status DESC, da.join_utc_time DESC, d.diary_id DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	err := r.db.Select(&diaries, baseSQL, args...)
	if err == sql.ErrNoRows {
		return []diary.Diary{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}

	count, err := r.GetJoinedDiariesCountByUserID(userID, includePrivate)
	if err != nil {
		return nil, 0, err
	}

	return diaries, count, nil
}

// GetJoinedDiariesCountByUserID 获取用户参与的日记本数量 (Active + Pending)
func (r *DiaryRepository) GetJoinedDiariesCountByUserID(userID int64, includePrivate bool) (int, error) {
	var count int
	baseSQL := `
		SELECT COUNT(*)
		FROM br_diary_author da
		JOIN br_diary d ON da.diary_id = d.diary_id
		WHERE da.user_id = ?`

	var args []interface{}
	args = append(args, userID)

	if includePrivate {
		baseSQL += " AND da.status IN (?, ?)"
		args = append(args, diary.AuthorStatusActive, diary.AuthorStatusPending)
	} else {
		baseSQL += " AND da.status = ?"
		args = append(args, diary.AuthorStatusActive)
	}

	baseSQL += " AND d.status = ? AND d.add_user_id != ?"
	args = append(args, diary.StatusActive, userID)

	if !includePrivate {
		baseSQL += " AND d.view_setting = ?"
		args = append(args, diary.ViewSettingPublic)
	}

	err := r.db.Get(&count, baseSQL, args...)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// GetActiveJoinedDiariesCountByUserID 获取用户参与的活跃日记本数量（不包含待确认）
func (r *DiaryRepository) GetActiveJoinedDiariesCountByUserID(userID int64) (int, error) {
	var count int
	sqlstr := `
		SELECT COUNT(*)
		FROM br_diary_author da
		JOIN br_diary d ON da.diary_id = d.diary_id
		WHERE da.user_id = ? AND da.status = ? AND d.status = ? AND d.add_user_id != ?
	`
	err := r.db.Get(&count, sqlstr, userID, diary.AuthorStatusActive, diary.StatusActive, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// GetDiaryAuthor 获取指定作者信息
func (r *DiaryRepository) GetDiaryAuthor(diaryID, userID int64) (*diary.Author, error) {
	var author diary.Author
	sqlstr := `SELECT * FROM br_diary_author WHERE diary_id = ? AND user_id = ?`
	err := r.db.Get(&author, sqlstr, diaryID, userID)
	if err == sql.ErrNoRows {
		return nil, errors.AuthorNotFound
	}
	return &author, err
}

// IsDiaryAuthor 检查用户是否为活跃作者
func (r *DiaryRepository) IsDiaryAuthor(diaryID, userID int64) (bool, error) {
	var status int8
	sqlstr := `SELECT status FROM br_diary_author WHERE diary_id = ? AND user_id = ?`
	err := r.db.Get(&status, sqlstr, diaryID, userID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return status == diary.AuthorStatusActive, nil
}

// ListPendingInvitationsByUserID 获取用户的待确认邀请列表
func (r *DiaryRepository) ListPendingInvitationsByUserID(userID int64) ([]diary.Diary, error) {
	var diaries []diary.Diary
	sqlstr := `
		SELECT d.*, u.handle, u.display_name
		FROM br_diary d
		JOIN br_diary_author da ON d.diary_id = da.diary_id
		LEFT JOIN br_user u ON d.add_user_id = u.user_id
		WHERE da.user_id = ? AND da.status = ? AND d.status = ?
		ORDER BY da.join_utc_time DESC
	`
	err := r.db.Select(&diaries, sqlstr, userID, diary.AuthorStatusPending, diary.StatusActive)
	if err == sql.ErrNoRows {
		return []diary.Diary{}, nil
	}
	return diaries, err
}

// IncrementAuthorPostNum 增加作者的发表数量
func (r *DiaryRepository) IncrementAuthorPostNum(diaryID, userID int64) error {
	sqlstr := "UPDATE br_diary_author SET post_num = post_num + 1 WHERE diary_id = ? AND user_id = ?"
	_, err := r.db.Exec(sqlstr, diaryID, userID)
	return err
}

// DecrementAuthorPostNum 减少作者的发表数量
func (r *DiaryRepository) DecrementAuthorPostNum(diaryID, userID int64) error {
	sqlstr := "UPDATE br_diary_author SET post_num = post_num - 1 WHERE diary_id = ? AND user_id = ? AND post_num > 0"
	_, err := r.db.Exec(sqlstr, diaryID, userID)
	return err
}
