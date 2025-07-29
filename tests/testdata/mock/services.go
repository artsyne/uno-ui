package mock

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/pkg"
)



// UserServiceMock 用户服务模拟器
type UserServiceMock struct {
	users     map[int64]*user.User
	Container *di.Container
}

// NewUserServiceMock 创建用户服务模拟器
func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{
		users: make(map[int64]*user.User),
	}
}

// 实现 user.ServiceInterface 接口
func (m *UserServiceMock) RegisterUser(u *user.User) error {
	// 验证输入参数
	if len(u.UserName) < 3 {
		return fmt.Errorf("输入参数错误")
	}
	if len(u.Pwd) < 6 {
		return fmt.Errorf("输入参数错误")
	}
	if u.TimeZone != "Asia/Shanghai" && u.TimeZone != "UTC" {
		return fmt.Errorf("输入参数错误")
	}

	// 生成用户ID
	if u.UserId <= 0 {
		u.UserId = int64(len(m.users) + 1)
	}

	// 使用 Bcrypt 加密密码
	originalPwd := u.Pwd
	hash, err := pkg.BcryptHash(originalPwd)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}
	u.Pwd = hash

	m.users[u.UserId] = u
	return nil
}

func (m *UserServiceMock) CheckNameAndPassword(name string, pwd string) (*user.User, error) {
	for _, u := range m.users {
		if u.UserName == name || u.Email == name {
			// 使用 Bcrypt 验证密码
			if pkg.BcryptVerify(pwd, u.Pwd) {
				return u, nil
			}
			return nil, fmt.Errorf("密码错误")
		}
	}
	return nil, fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) ActivateUser(uid string, sign string) bool {
	return true
}

func (m *UserServiceMock) IsExistEmail(email string) error {
	for _, u := range m.users {
		if u.Email == email {
			return fmt.Errorf("邮箱已存在")
		}
	}
	return nil
}

func (m *UserServiceMock) IsExistUserName(userName string) error {
	for _, u := range m.users {
		if u.UserName == userName {
			return fmt.Errorf("用户名已存在")
		}
	}
	return nil
}

func (m *UserServiceMock) CreateResetPasswordCode(userID int64) error {
	return nil
}

func (m *UserServiceMock) SendResetPasswordEmail(u *user.User) error {
	return nil
}

func (m *UserServiceMock) ChangeUserPassword(currentUserID int64, currentPwd, newPwd string) error {
	// 检查用户是否存在
	user, exists := m.users[currentUserID]
	if !exists {
		return fmt.Errorf("用户不存在")
	}

	// 检查新密码长度
	if len(newPwd) < 6 {
		return fmt.Errorf("新密码长度不能少于6位")
	}

	// 使用 Bcrypt 验证当前密码
	if !pkg.BcryptVerify(currentPwd, user.Pwd) {
		return fmt.Errorf("当前密码错误")
	}

	// 使用 Bcrypt 加密新密码
	hash, err := pkg.BcryptHash(newPwd)
	if err != nil {
		return fmt.Errorf("新密码加密失败: %v", err)
	}
	user.Pwd = hash
	return nil
}

func (m *UserServiceMock) ResetUserPassword(userId int64, newPassword string) error {
	return nil
}

func (m *UserServiceMock) CheckValidCode(userID int64, sign string, validType int8) (bool, int64) {
	return true, 1
}

func (m *UserServiceMock) DeleteValidCode(id int64) (int64, error) {
	return 1, nil
}

func (m *UserServiceMock) GetUserByEmail(email string) (*user.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) GetUserByUserId(userId int64) (*user.User, error) {
	if u, exists := m.users[userId]; exists {
		return u, nil
	}
	return nil, fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) ParseUserID(idStr string) (int64, error) {
	return 1, nil
}

func (m *UserServiceMock) GetUserProfile(uid int64, currentUserID int64) (*user.User, bool, error) {
	if u, exists := m.users[uid]; exists {
		return u, true, nil
	}
	return nil, false, fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) GetUserSettings(uid int64) (*user.User, error) {
	if u, exists := m.users[uid]; exists {
		return u, nil
	}
	return nil, fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) UpdateUserProfile(uid int64, currentUserID int64, fullName, remark string) error {
	if u, exists := m.users[uid]; exists {
		u.FullName = fullName
		u.Remark = remark
		return nil
	}
	return fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) UpdateUserEmail(uid int64, currentUserID int64, email string) error {
	if u, exists := m.users[uid]; exists {
		u.Email = email
		return nil
	}
	return fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) UpdateUserStatus(userId int64, status int8) (int64, error) {
	if u, exists := m.users[userId]; exists {
		u.Status = status
		return 1, nil
	}
	return 0, fmt.Errorf("用户不存在")
}

func (m *UserServiceMock) UpdateUserPwd(userId int64, pwd string) (int64, error) {
	if user, exists := m.users[userId]; exists {
		user.Pwd = pwd
		return 1, nil
	}
	return 0, fmt.Errorf("用户不存在")
}

// GetSupportedLang 获取支持的语言
func (m *UserServiceMock) GetSupportedLang(language string) string {
	supportedLangs := []string{"zh-cn", "zh-tw", "en"}
	for _, lang := range supportedLangs {
		if lang == language {
			return language
		}
	}
	return "zh-cn" // 默认中文
}

// SetupMockUserService 设置模拟用户服务
func SetupMockUserService(t *testing.T, container *di.Container) *UserServiceMock {
	t.Helper()

	mock := NewUserServiceMock()

	// 添加测试用户
	testUser := &user.User{
		UserId:   1,
		UserName: "testuser",
		Email:    "test@example.com",
		Status:   1,
		FullName: "Test User",
		Lang:     "en",
		TimeZone: "UTC",
		Pwd:      "encrypted_password123", // 添加默认密码（加密格式）
	}
	mock.users[1] = testUser

	mock.Container = container

	// 绑定用户服务到DI容器
	err := di.Bind(container, func() user.ServiceInterface {
		return mock
	})
	if err != nil {
		t.Fatalf("绑定用户服务到DI容器失败: %v", err)
	}

	return mock
}



// DiaryServiceMock 日记服务模拟器
type DiaryServiceMock struct {
	diaries   map[int64]*diary.Diary
	Container *di.Container
}

// NewDiaryServiceMock 创建日记服务模拟器
func NewDiaryServiceMock() *DiaryServiceMock {
	return &DiaryServiceMock{
		diaries: make(map[int64]*diary.Diary),
	}
}

// 实现 diary.ServiceInterface 接口
func (m *DiaryServiceMock) AddDiary(d *diary.Diary) (int64, error) {
	id := int64(len(m.diaries) + 1)
	d.DiaryID = id
	m.diaries[id] = d
	return id, nil
}

func (m *DiaryServiceMock) UpdateDiary(d *diary.Diary) error {
	if _, exists := m.diaries[d.DiaryID]; exists {
		m.diaries[d.DiaryID] = d
		return nil
	}
	return fmt.Errorf("日记不存在")
}

func (m *DiaryServiceMock) GetDiaryByID(diaryId int64) (*diary.Diary, error) {
	if d, exists := m.diaries[diaryId]; exists {
		return d, nil
	}
	return nil, fmt.Errorf("日记不存在")
}

func (m *DiaryServiceMock) ListDiariesForWrite(userId int64) ([]diary.Diary, error) {
	var result []diary.Diary
	for _, d := range m.diaries {
		if d.AddUserID == userId {
			result = append(result, *d)
		}
	}
	return result, nil
}

func (m *DiaryServiceMock) ListDiariesByUser(userId int64, page, pageSize int) ([]diary.Diary, int, error) {
	var result []diary.Diary
	count := 0
	for _, d := range m.diaries {
		if d.AddUserID == userId {
			result = append(result, *d)
			count++
		}
	}
	return result, count, nil
}

func (m *DiaryServiceMock) IncrementPostCount(diaryId int64, userId int64, userName string, actionTime int64) error {
	if d, exists := m.diaries[diaryId]; exists {
		d.InfoNum++
		d.LastActionUserID = userId
		d.LastActionUserName = userName
		d.LastActionUTCTime = actionTime
		return nil
	}
	return fmt.Errorf("日记不存在")
}

func (m *DiaryServiceMock) DecrementPostCount(diaryId int64, userId int64, userName string, actionTime int64) error {
	if d, exists := m.diaries[diaryId]; exists {
		if d.InfoNum > 0 {
			d.InfoNum--
		}
		d.LastActionUserID = userId
		d.LastActionUserName = userName
		d.LastActionUTCTime = actionTime
		return nil
	}
	return fmt.Errorf("日记不存在")
}

func (m *DiaryServiceMock) CheckDiaryLimit(userId int64) error {
	count := 0
	for _, d := range m.diaries {
		if d.AddUserID == userId {
			count++
		}
	}
	if count >= 10 { // 假设限制为10个日记本
		return fmt.Errorf("日记本数量已达上限")
	}
	return nil
}

// SetupMockDiaryService 设置模拟日记服务
func SetupMockDiaryService(t *testing.T, container *di.Container) *DiaryServiceMock {
	t.Helper()

	mock := NewDiaryServiceMock()
	mock.Container = container

	// 绑定日记服务到DI容器
	err := di.Bind(container, func() diary.ServiceInterface {
		return mock
	})
	if err != nil {
		t.Fatalf("绑定日记服务到DI容器失败: %v", err)
	}

	return mock
}



// PostServiceMock 文章服务模拟器
type PostServiceMock struct {
	posts     map[int64]*post.Post
	Container *di.Container
}

// NewPostServiceMock 创建文章服务模拟器
func NewPostServiceMock() *PostServiceMock {
	return &PostServiceMock{
		posts: make(map[int64]*post.Post),
	}
}

// 实现 post.ServiceInterface 接口
func (m *PostServiceMock) AddPost(p *post.Post) (int64, error) {
	id := int64(len(m.posts) + 1)
	p.PostID = id
	m.posts[id] = p
	return id, nil
}

func (m *PostServiceMock) GetPostByID(postId int64) (*post.Post, error) {
	if p, exists := m.posts[postId]; exists {
		return p, nil
	}
	return nil, fmt.Errorf("文章不存在")
}

func (m *PostServiceMock) UpdatePost(p *post.Post) error {
	if _, exists := m.posts[p.PostID]; exists {
		m.posts[p.PostID] = p
		return nil
	}
	return fmt.Errorf("文章不存在")
}

func (m *PostServiceMock) ListPostsByUserID(userId int64, diaryId int64, page, pageSize int) ([]post.Post, int, error) {
	var result []post.Post
	count := 0
	for _, p := range m.posts {
		if p.AddUserID == userId && p.DiaryID == diaryId {
			result = append(result, *p)
			count++
		}
	}
	return result, count, nil
}

func (m *PostServiceMock) GetPostCountByUserID(userId int64, diaryId int64) (int, error) {
	count := 0
	for _, p := range m.posts {
		if p.AddUserID == userId && p.DiaryID == diaryId {
			count++
		}
	}
	return count, nil
}

func (m *PostServiceMock) UpdatePostByID(p *post.Post) error {
	if _, exists := m.posts[p.PostID]; exists {
		m.posts[p.PostID] = p
		return nil
	}
	return fmt.Errorf("文章不存在")
}

func (m *PostServiceMock) DeletePostByID(postID int64) error {
	if _, exists := m.posts[postID]; exists {
		delete(m.posts, postID)
		return nil
	}
	return fmt.Errorf("文章不存在")
}

// SetupMockPostService 设置模拟文章服务
func SetupMockPostService(t *testing.T, container *di.Container) *PostServiceMock {
	t.Helper()

	mock := NewPostServiceMock()
	mock.Container = container

	// 绑定文章服务到DI容器
	err := di.Bind(container, func() post.ServiceInterface {
		return mock
	})
	if err != nil {
		t.Fatalf("绑定文章服务到DI容器失败: %v", err)
	}

	return mock
}



// TaskServiceMock 任务服务模拟器
type TaskServiceMock struct {
	tasks     map[int64]*task.EmailTask
	Container *di.Container
}

// NewTaskServiceMock 创建任务服务模拟器
func NewTaskServiceMock() *TaskServiceMock {
	return &TaskServiceMock{
		tasks: make(map[int64]*task.EmailTask),
	}
}

// 实现 task.ServiceInterface 接口
func (m *TaskServiceMock) StartScheduler() {
	// Mock implementation - do nothing
}

func (m *TaskServiceMock) ProcessEmailTasks() {
	// Mock implementation - do nothing
}

func (m *TaskServiceMock) CreateEmailTask(to, subject, body string) error {
	id := int64(len(m.tasks) + 1)
	emailTask := &task.EmailTask{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	m.tasks[id] = emailTask
	return nil
}

// SetupMockTaskService 设置模拟任务服务
func SetupMockTaskService(t *testing.T, container *di.Container) *TaskServiceMock {
	t.Helper()

	mock := NewTaskServiceMock()
	mock.Container = container

	// 绑定任务服务到DI容器
	err := di.Bind(container, func() task.ServiceInterface {
		return mock
	})
	if err != nil {
		t.Fatalf("绑定任务服务到DI容器失败: %v", err)
	}

	return mock
}
