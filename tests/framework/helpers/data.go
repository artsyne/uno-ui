package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
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
	rand.Read(randomBytes)

	return &DataGeneratorImpl{
		isolator: isolator,
		random:   hex.EncodeToString(randomBytes),
	}
}

// CreateUser 创建测试用户
func (dg *DataGeneratorImpl) CreateUser(options ...UserOption) *user.User {
	// 默认选项
	opts := &UserOptions{
		Username: dg.isolator.IsolateUsername("user"),
		Email:    dg.isolator.IsolateEmail("user@test.com"),
		Password: "password123",
		FullName: "Test User",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
	}

	// 应用自定义选项
	for _, opt := range options {
		opt(opts)
	}

	return &user.User{
		UserName:         opts.Username,
		Email:            opts.Email,
		Pwd:              opts.Password,
		FullName:         opts.FullName,
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
		Title:             opts.Title,
		Remark:            opts.Description,
		ViewSetting:       int8(opts.ViewSetting),
		SocialSetting:     int8(opts.SocialSetting),
		ThemeID:           int8(opts.ThemeID),
		AddUTCTime:        now,
		AddUserID:         userID,
		LastActionUTCTime: now,
		LastActionUserID:  userID,
		Status:            1, // 正常状态
	}
}

// CreatePost 创建测试文章
func (dg *DataGeneratorImpl) CreatePost(userID, diaryID int64, options ...PostOption) *post.Post {
	// 默认选项
	opts := &PostOptions{
		Title:    dg.isolator.IsolateTitle("post"),
		Body:     "Test post body content with detailed information",
		Summary:  "Test post summary",
		InfoType: 1, // 普通文章
	}

	// 应用自定义选项
	for _, opt := range options {
		opt(opts)
	}

	now := time.Now().Unix()

	return &post.Post{
		Title:         opts.Title,
		Body:          opts.Body,
		Summary:       opts.Summary,
		DiaryID:       diaryID,
		AddUserID:     userID,
		AddUTCTime:    now,
		UpdateUTCTime: now,
		Status:        1, // 正常状态
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

// GenerateUniqueUsername 生成唯一用户名
func (dg *DataGeneratorImpl) GenerateUniqueUsername(prefix string) string {
	return dg.isolator.IsolateUsername(prefix)
}

// GenerateUniqueTitle 生成唯一标题
func (dg *DataGeneratorImpl) GenerateUniqueTitle(prefix string) string {
	return dg.isolator.IsolateTitle(prefix)
}

// User选项函数

// WithUsername 设置用户名
func WithUsername(username string) UserOption {
	return func(opts *UserOptions) {
		opts.Username = username
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

// WithFullName 设置全名
func WithFullName(fullName string) UserOption {
	return func(opts *UserOptions) {
		opts.FullName = fullName
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

// Post选项函数

// WithPostTitle 设置文章标题
func WithPostTitle(title string) PostOption {
	return func(opts *PostOptions) {
		opts.Title = title
	}
}

// WithPostBody 设置文章内容
func WithPostBody(body string) PostOption {
	return func(opts *PostOptions) {
		opts.Body = body
	}
}

// WithPostSummary 设置文章摘要
func WithPostSummary(summary string) PostOption {
	return func(opts *PostOptions) {
		opts.Summary = summary
	}
}

// WithInfoType 设置信息类型
func WithInfoType(infoType int) PostOption {
	return func(opts *PostOptions) {
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
			WithUsername(fmt.Sprintf("user%d", i)),
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

// CreatePosts 批量创建文章
func (dg *DataGeneratorImpl) CreatePosts(userID, diaryID int64, count int, baseOptions ...PostOption) []*post.Post {
	posts := make([]*post.Post, count)

	for i := 0; i < count; i++ {
		// 为每个文章创建唯一的选项
		options := make([]PostOption, len(baseOptions))
		copy(options, baseOptions)

		// 添加索引以确保唯一性
		options = append(options,
			WithPostTitle(fmt.Sprintf("post%d", i)),
		)

		posts[i] = dg.CreatePost(userID, diaryID, options...)
	}

	return posts
}

// CreateCompleteUserScenario 创建完整的用户场景（用户+日记+文章）
func (dg *DataGeneratorImpl) CreateCompleteUserScenario(diaryCount, postCount int) (*user.User, []*diary.Diary, []*post.Post) {
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
	posts := make([]*post.Post, 0)
	for i := range diaries {
		// 每个日记创建指定数量的文章
		diaryID := int64(2000 + i)
		diaryPosts := dg.CreatePosts(userID, diaryID, postCount,
			WithPostTitle(fmt.Sprintf("diary%d_post", i+1)),
		)
		posts = append(posts, diaryPosts...)
	}

	return testUser, diaries, posts
}
