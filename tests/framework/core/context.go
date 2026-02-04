package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/user"
	. "github.com/biuren/hebe/tests/framework/types"
)

// DataGeneratorImpl 测试数据生成器实现
// 统一的测试数据创建，替代分散的数据生成逻辑
type DataGeneratorImpl struct {
	isolator DataIsolatorInterface
	random   string
}

// NewDataGenerator 创建数据生成器
func NewDataGenerator(isolator DataIsolatorInterface) *DataGeneratorImpl {
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}

	return &DataGeneratorImpl{
		isolator: isolator,
		random:   hex.EncodeToString(randomBytes),
	}
}

// CreateUser 创建测试用户
func (dg *DataGeneratorImpl) CreateUser(options ...UserOption) *user.User {
	// 默认选项
	opts := &UserOptions{
		Handle:      dg.isolator.IsolateHandle("user"),
		Email:       dg.isolator.IsolateEmail("user@test.com"),
		Password:    "password123",
		DisplayName: "Test User",
		TimeZone:    "Asia/Shanghai",
		Lang:        "zh-cn",
	}

	// 应用自定义选项
	for _, opt := range options {
		opt(opts)
	}

	return &user.User{
		Handle:           opts.Handle,
		Email:            opts.Email,
		Pwd:              opts.Password,
		DisplayName:      opts.DisplayName,
		TimeZone:         opts.TimeZone,
		Lang:             opts.Lang,
		CreatedUTCTime:   time.Now().Unix(),
		LastLoginUTCTime: time.Now().Unix(),
		Status:           3, // 正常状态：3
	}
}

// CreateDiary 创建测试日记本
func (dg *DataGeneratorImpl) CreateDiary(userID int64, options ...DiaryOption) *diary.Diary {
	// 默认选项
	opts := &DiaryOptions{
		Title:         dg.isolator.IsolateTitle("diary"),
		Description:   "Test diary description",
		ViewSetting:   1, // 公开
		SocialSetting: 1, // 允许社交
		ThemeID:       1, // 默认主题
	}

	// 应用自定义选项
	for _, opt := range options {
		opt(opts)
	}

	now := time.Now().Unix()

	return &diary.Diary{
		Title:         opts.Title,
		Remark:        opts.Description,
		ViewSetting:   int8(opts.ViewSetting),
		SocialSetting: int8(opts.SocialSetting),
		ThemeID:       int8(opts.ThemeID),
		AddUTCTime:    now,
		AddUserID:     userID,
		Status:        1, // 正常状态
	}
}

// CreateEntry 创建测试文章
func (dg *DataGeneratorImpl) CreateEntry(userID, diaryID int64, options ...EntryOption) *entry.Entry {
	// 默认选项
	opts := &EntryOptions{
		Title:    dg.isolator.IsolateTitle("entry"),
		Body:     "Test entry body content with detailed information",
		Summary:  "Test entry summary",
		InfoType: entry.InfoTypeText, // 普通文章
	}

	// 应用自定义选项
	for _, opt := range options {
		opt(opts)
	}

	now := time.Now().Unix()

	return &entry.Entry{
		DiaryID:       diaryID,
		Title:         opts.Title,
		Summary:       opts.Summary,
		Body:          opts.Body,
		InfoType:      opts.InfoType,
		Status:        1, // 正常状态
		AddUTCTime:    now,
		UpdateUTCTime: now,
		AddUserID:     userID,
		AddUserIP:     "127.0.0.1",
		ViewNum:       0,
		LikeNum:       0,
		CommentNum:    0,
	}
}

// GenerateUniqueString 生成唯一字符串
func (dg *DataGeneratorImpl) GenerateUniqueString(prefix string) string {
	return dg.isolator.IsolateString(prefix)
}

// GenerateUniqueEmail 生成唯一邮箱
func (dg *DataGeneratorImpl) GenerateUniqueEmail(prefix string) string {
	email := fmt.Sprintf("%s@test.com", prefix)
	return dg.isolator.IsolateEmail(email)
}

// GenerateUniqueHandle 生成唯一Handle
func (dg *DataGeneratorImpl) GenerateUniqueHandle(prefix string) string {
	return dg.isolator.IsolateHandle(prefix)
}

// GenerateUniqueTitle 生成唯一标题
func (dg *DataGeneratorImpl) GenerateUniqueTitle(prefix string) string {
	return dg.isolator.IsolateTitle(prefix)
}

// User选项函数

// WithHandle 设置Handle
func WithHandle(handle string) UserOption {
	return func(opts *UserOptions) {
		opts.Handle = handle
	}
}

// WithEmail 设置邮箱
func WithEmail(email string) UserOption {
	return func(opts *UserOptions) {
		opts.Email = email
	}
}

// WithPassword 设置密码
func WithPassword(password string) UserOption {
	return func(opts *UserOptions) {
		opts.Password = password
	}
}

// WithDisplayName 设置显示名称
func WithDisplayName(displayName string) UserOption {
	return func(opts *UserOptions) {
		opts.DisplayName = displayName
	}
}

// WithTimeZone 设置时区
func WithTimeZone(timeZone string) UserOption {
	return func(opts *UserOptions) {
		opts.TimeZone = timeZone
	}
}

// WithLang 设置语言
func WithLang(lang string) UserOption {
	return func(opts *UserOptions) {
		opts.Lang = lang
	}
}

// Diary选项函数

// WithDiaryTitle 设置日记标题
func WithDiaryTitle(title string) DiaryOption {
	return func(opts *DiaryOptions) {
		opts.Title = title
	}
}

// WithDiaryDescription 设置日记描述
func WithDiaryDescription(description string) DiaryOption {
	return func(opts *DiaryOptions) {
		opts.Description = description
	}
}

// WithViewSetting 设置查看权限
func WithViewSetting(setting int) DiaryOption {
	return func(opts *DiaryOptions) {
		opts.ViewSetting = setting
	}
}

// WithSocialSetting 设置社交权限
func WithSocialSetting(setting int) DiaryOption {
	return func(opts *DiaryOptions) {
		opts.SocialSetting = setting
	}
}

// WithThemeID 设置主题ID
func WithThemeID(themeID int) DiaryOption {
	return func(opts *DiaryOptions) {
		opts.ThemeID = themeID
	}
}

// Entry选项函数

// WithEntryTitle 设置文章标题
func WithEntryTitle(title string) EntryOption {
	return func(opts *EntryOptions) {
		opts.Title = title
	}
}

// WithEntryBody 设置文章内容
func WithEntryBody(body string) EntryOption {
	return func(opts *EntryOptions) {
		opts.Body = body
	}
}

// WithEntrySummary 设置文章摘要
func WithEntrySummary(summary string) EntryOption {
	return func(opts *EntryOptions) {
		opts.Summary = summary
	}
}

// WithInfoType 设置信息类型
func WithInfoType(infoType entry.InfoType) EntryOption {
	return func(opts *EntryOptions) {
		opts.InfoType = infoType
	}
}

// 批量创建方法

// CreateUsers 批量创建用户
func (dg *DataGeneratorImpl) CreateUsers(count int, baseOptions ...UserOption) []*user.User {
	users := make([]*user.User, count)

	for i := 0; i < count; i++ {
		// 为每个用户创建唯一的选项
		options := make([]UserOption, len(baseOptions))
		copy(options, baseOptions)

		// 添加索引以确保唯一性
		options = append(options,
			WithHandle(fmt.Sprintf("user%d", i)),
			WithEmail(fmt.Sprintf("user%d@test.com", i)),
		)

		users[i] = dg.CreateUser(options...)
	}

	return users
}

// CreateDiaries 批量创建日记本
func (dg *DataGeneratorImpl) CreateDiaries(userID int64, count int, baseOptions ...DiaryOption) []*diary.Diary {
	diaries := make([]*diary.Diary, count)

	for i := 0; i < count; i++ {
		// 为每个日记创建唯一的选项
		options := make([]DiaryOption, len(baseOptions))
		copy(options, baseOptions)

		// 添加索引以确保唯一性
		options = append(options,
			WithDiaryTitle(fmt.Sprintf("diary%d", i)),
		)

		diaries[i] = dg.CreateDiary(userID, options...)
	}

	return diaries
}

// CreateEntries 批量创建文章
func (dg *DataGeneratorImpl) CreateEntries(userID, diaryID int64, count int, baseOptions ...EntryOption) []*entry.Entry {
	entries := make([]*entry.Entry, count)

	for i := 0; i < count; i++ {
		// 为每个文章创建唯一的选项
		options := make([]EntryOption, len(baseOptions))
		copy(options, baseOptions)

		// 添加索引以确保唯一性
		options = append(options,
			WithEntryTitle(fmt.Sprintf("entry%d", i)),
		)

		entries[i] = dg.CreateEntry(userID, diaryID, options...)
	}

	return entries
}

// CreateCompleteUserScenario 创建完整的用户场景（用户+日记+文章）
func (dg *DataGeneratorImpl) CreateCompleteUserScenario(diaryCount, entryCount int) (*user.User, []*diary.Diary, []*entry.Entry) {
	// 创建用户
	testUser := dg.CreateUser()

	// 模拟用户ID（在实际使用中会从数据库获取）
	userID := int64(1000 + time.Now().Unix()%1000)

	// 创建日记本
	diaries := make([]*diary.Diary, diaryCount)
	for i := 0; i < diaryCount; i++ {
		diaries[i] = dg.CreateDiary(userID,
			WithDiaryTitle(fmt.Sprintf("diary%d", i+1)),
		)
	}

	// 创建文章
	entries := make([]*entry.Entry, 0)
	for i := range diaries {
		// 每个日记创建指定数量的文章
		diaryID := int64(2000 + i)
		diaryEntries := dg.CreateEntries(userID, diaryID, entryCount,
			WithEntryTitle(fmt.Sprintf("diary%d_entry", i+1)),
		)
		entries = append(entries, diaryEntries...)
	}

	return testUser, diaries, entries
}
