# 代码质量改进指南

## 目录
1. [Go 后端代码规范](#go-后端代码规范)
2. [前端代码规范](#前端代码规范)
3. [测试规范](#测试规范)
4. [性能优化实践](#性能优化实践)
5. [安全编码规范](#安全编码规范)

## Go 后端代码规范

### 1. 包和文件组织

#### 问题示例
```go
// ❌ 错误：文件名不规范
pkg/Paginator.go  // 应该是 paginator.go

// ❌ 错误：包引入混乱
import (
    "github.com/biuren/hebe/pkg"  // 不具体
)
```

#### 最佳实践
```go
// ✅ 正确：清晰的包结构
pkg/
├── pagination/
│   └── paginator.go
├── crypto/
│   ├── hash.go
│   └── encrypt.go
└── validation/
    └── validator.go

// ✅ 正确：具体的包引入
import (
    "github.com/biuren/hebe/pkg/pagination"
    "github.com/biuren/hebe/pkg/crypto"
)
```

### 2. 接口设计

#### 问题示例
```go
// ❌ 错误：接口过大
type UserRepository interface {
    Create(user *User) error
    Update(user *User) error
    Delete(id int) error
    FindByID(id int) (*User, error)
    FindByEmail(email string) (*User, error)
    FindByUsername(username string) (*User, error)
    FindAll() ([]*User, error)
    FindWithPagination(page, size int) ([]*User, error)
    Count() (int, error)
    // ... 更多方法
}
```

#### 最佳实践
```go
// ✅ 正确：接口分离原则
type UserReader interface {
    FindByID(id int) (*User, error)
    FindByEmail(email string) (*User, error)
}

type UserWriter interface {
    Create(user *User) error
    Update(user *User) error
    Delete(id int) error
}

type UserRepository interface {
    UserReader
    UserWriter
}

// ✅ 正确：按需定义小接口
type Authenticator interface {
    Authenticate(email, password string) (*User, error)
}
```

### 3. 错误处理

#### 问题示例
```go
// ❌ 错误：错误信息不明确
if err != nil {
    return errors.New("error")
}

// ❌ 错误：忽略错误
result, _ := someFunction()

// ❌ 错误：错误处理不一致
if user == nil {
    return nil, fmt.Errorf("user not found")  // 应使用 AppError
}
```

#### 最佳实践
```go
// ✅ 正确：使用预定义错误
if err != nil {
    return errors.DatabaseError.WithDetails(err.Error())
}

// ✅ 正确：始终处理错误
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// ✅ 正确：错误包装
if user == nil {
    return nil, errors.UserNotFound.WithDetails(
        fmt.Sprintf("user id: %d", userID),
    )
}

// ✅ 正确：错误日志记录
func (s *Service) CreateUser(user *User) error {
    if err := s.repo.Create(user); err != nil {
        slog.Error("failed to create user",
            "error", err,
            "email", user.Email,
            "trace_id", ctx.Value("trace_id"),
        )
        return errors.DatabaseError.WithDetails("create user failed")
    }
    return nil
}
```

### 4. 并发安全

#### 问题示例
```go
// ❌ 错误：数据竞争
type Counter struct {
    count int
}

func (c *Counter) Inc() {
    c.count++  // 非线程安全
}
```

#### 最佳实践
```go
// ✅ 正确：使用互斥锁
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

// ✅ 正确：使用原子操作
type AtomicCounter struct {
    count int64
}

func (c *AtomicCounter) Inc() {
    atomic.AddInt64(&c.count, 1)
}

// ✅ 正确：使用 channel
type SafeCounter struct {
    ch chan int
}

func NewSafeCounter() *SafeCounter {
    c := &SafeCounter{
        ch: make(chan int, 1),
    }
    c.ch <- 0
    return c
}

func (c *SafeCounter) Inc() {
    count := <-c.ch
    c.ch <- count + 1
}
```

### 5. Context 使用

#### 问题示例
```go
// ❌ 错误：没有传递 context
func (s *Service) GetUser(id int) (*User, error) {
    return s.repo.FindByID(id)
}
```

#### 最佳实践
```go
// ✅ 正确：传递 context
func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
    // 检查 context 是否已取消
    if err := ctx.Err(); err != nil {
        return nil, err
    }

    // 设置操作超时
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return s.repo.FindByID(ctx, id)
}

// ✅ 正确：在 context 中传递请求信息
func middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        ctx = context.WithValue(ctx, "request_id", generateRequestID())
        ctx = context.WithValue(ctx, "user_id", getUserID(r))

        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
```

### 6. 资源管理

#### 问题示例
```go
// ❌ 错误：资源泄露
func readFile(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    // 忘记关闭文件
    return io.ReadAll(file)
}
```

#### 最佳实践
```go
// ✅ 正确：使用 defer 确保资源释放
func readFile(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    return io.ReadAll(file)
}

// ✅ 正确：数据库连接管理
func (s *Service) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()

    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}
```

## 前端代码规范

### 1. JavaScript 模块化

#### 问题示例
```javascript
// ❌ 错误：全局变量污染
var userList = [];
function loadUsers() {
    // ...
}

// ❌ 错误：没有模块化
window.app = window.app || {};
window.app.users = {
    load: function() {}
};
```

#### 最佳实践
```javascript
// ✅ 正确：ES6 模块
// users.js
export class UserManager {
    constructor(apiClient) {
        this.apiClient = apiClient;
        this.users = [];
    }

    async load() {
        this.users = await this.apiClient.get('/api/users');
        return this.users;
    }
}

// ✅ 正确：Alpine.js 组件化
Alpine.data('userList', () => ({
    users: [],
    loading: false,

    async init() {
        this.loading = true;
        try {
            this.users = await this.$store.api.get('/api/users');
        } catch (error) {
            this.$store.notify.error('Failed to load users');
        } finally {
            this.loading = false;
        }
    }
}));
```

### 2. 异步处理

#### 问题示例
```javascript
// ❌ 错误：回调地狱
getData(function(a) {
    getMoreData(a, function(b) {
        getMoreData(b, function(c) {
            console.log(c);
        });
    });
});

// ❌ 错误：未处理错误
async function loadData() {
    const data = await fetch('/api/data');  // 可能失败
    return data.json();
}
```

#### 最佳实践
```javascript
// ✅ 正确：使用 async/await 和错误处理
async function loadData() {
    try {
        const response = await fetch('/api/data');

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.json();
    } catch (error) {
        console.error('Failed to load data:', error);
        // 通知用户
        this.$store.notify.error('数据加载失败，请稍后重试');
        return null;
    }
}

// ✅ 正确：并发请求优化
async function loadDashboard() {
    // 并发执行多个请求
    const [users, stats, activities] = await Promise.all([
        this.loadUsers(),
        this.loadStats(),
        this.loadActivities()
    ]);

    return { users, stats, activities };
}
```

### 3. 状态管理

#### 问题示例
```javascript
// ❌ 错误：状态分散
let currentUser = null;
let isLoggedIn = false;
let userPermissions = [];

// ❌ 错误：直接修改 DOM
document.getElementById('username').innerText = user.name;
```

#### 最佳实践
```javascript
// ✅ 正确：集中状态管理 (Alpine.js Store)
Alpine.store('auth', {
    user: null,
    permissions: [],

    get isLoggedIn() {
        return this.user !== null;
    },

    login(user) {
        this.user = user;
        this.permissions = user.permissions || [];
    },

    logout() {
        this.user = null;
        this.permissions = [];
    },

    hasPermission(permission) {
        return this.permissions.includes(permission);
    }
});

// ✅ 正确：响应式数据绑定
<div x-data>
    <span x-text="$store.auth.user?.name || '游客'"></span>
</div>
```

## 测试规范

### 1. 单元测试

#### 最佳实践
```go
// ✅ 正确：表格驱动测试
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {
            name:    "valid email",
            email:   "user@example.com",
            wantErr: false,
        },
        {
            name:    "invalid format",
            email:   "invalid-email",
            wantErr: true,
        },
        {
            name:    "empty email",
            email:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// ✅ 正确：Mock 接口
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}
```

### 2. 集成测试

#### 最佳实践
```go
// ✅ 正确：测试隔离
func TestUserAPI(t *testing.T) {
    // 创建测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(db)

    // 创建测试服务器
    server := setupTestServer(db)

    // 测试用例
    t.Run("CreateUser", func(t *testing.T) {
        user := &User{
            Email:    "test@example.com",
            Password: "password123",
        }

        resp := httptest.NewRecorder()
        req := httptest.NewRequest("POST", "/api/users", toJSON(user))

        server.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusCreated, resp.Code)
    })
}
```

## 性能优化实践

### 1. 数据库优化

#### 最佳实践
```go
// ✅ 正确：批量操作
func (r *Repository) CreateBatch(users []*User) error {
    query := `INSERT INTO users (email, name, created_at) VALUES `
    values := []interface{}{}

    for i, user := range users {
        if i > 0 {
            query += ","
        }
        query += "(?, ?, ?)"
        values = append(values, user.Email, user.Name, time.Now())
    }

    _, err := r.db.Exec(query, values...)
    return err
}

// ✅ 正确：预加载关联
func (r *Repository) FindDiariesWithAuthors() ([]*Diary, error) {
    query := `
        SELECT
            d.id, d.title, d.content, d.author_id,
            u.id, u.name, u.email
        FROM diaries d
        LEFT JOIN users u ON d.author_id = u.id
        ORDER BY d.created_at DESC
    `

    // 使用 sqlx 的结构体扫描
    var results []struct {
        Diary
        Author User `db:"u"`
    }

    err := r.db.Select(&results, query)
    // 处理结果...
}
```

### 2. 缓存策略

#### 最佳实践
```go
// ✅ 正确：多级缓存
type CacheService struct {
    local *lru.Cache     // 本地缓存
    redis *redis.Client  // 分布式缓存
}

func (c *CacheService) Get(key string) (interface{}, error) {
    // 1. 检查本地缓存
    if val, ok := c.local.Get(key); ok {
        return val, nil
    }

    // 2. 检查 Redis
    val, err := c.redis.Get(key).Result()
    if err == nil {
        c.local.Add(key, val)  // 更新本地缓存
        return val, nil
    }

    return nil, err
}

// ✅ 正确：缓存预热
func (s *Service) PrewarmCache(ctx context.Context) error {
    // 预加载热点数据
    hotUsers, err := s.repo.FindHotUsers(ctx, 100)
    if err != nil {
        return err
    }

    for _, user := range hotUsers {
        key := fmt.Sprintf("user:%d", user.ID)
        s.cache.Set(key, user, 1*time.Hour)
    }

    return nil
}
```

## 安全编码规范

### 1. 输入验证

#### 最佳实践
```go
// ✅ 正确：严格的输入验证
func ValidateUserInput(input *UserInput) error {
    // 长度检查
    if len(input.Username) < 3 || len(input.Username) > 20 {
        return errors.InvalidUsername
    }

    // 格式检查
    if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(input.Username) {
        return errors.InvalidUsername
    }

    // SQL 注入防护
    if containsSQLKeywords(input.Username) {
        return errors.InvalidInput
    }

    // XSS 防护
    input.Bio = html.EscapeString(input.Bio)

    return nil
}

// ✅ 正确：参数化查询
func (r *Repository) FindByUsername(username string) (*User, error) {
    query := "SELECT * FROM users WHERE username = ? LIMIT 1"
    user := &User{}
    err := r.db.Get(user, query, username)
    return user, err
}
```

### 2. 认证和授权

#### 最佳实践
```go
// ✅ 正确：JWT 令牌管理
type TokenService struct {
    secret []byte
    issuer string
}

func (s *TokenService) GenerateToken(user *User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
        "iss":     s.issuer,
        "iat":     time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secret)
}

// ✅ 正确：中间件权限检查
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(*User)

        if !user.HasPermission(permission) {
            c.JSON(403, gin.H{
                "error": "Permission denied",
                "code":  1201,
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## 总结

遵循这些代码质量规范可以：

1. **提高代码可读性** - 统一的编码风格让团队协作更顺畅
2. **减少 Bug** - 正确的错误处理和资源管理避免常见问题
3. **提升性能** - 优化的查询和缓存策略提高响应速度
4. **增强安全性** - 严格的输入验证和认证机制保护系统
5. **易于维护** - 模块化设计和充分的测试让重构更安全

建议定期进行代码审查，确保这些规范得到遵守。