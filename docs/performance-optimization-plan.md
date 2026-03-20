# 性能优化实施方案

## 1. 数据库性能优化

### 1.1 查询优化

#### 现状分析
通过代码审查发现以下问题：
- 存在 N+1 查询问题
- 缺少适当的索引
- 部分查询未使用分页

#### 优化方案

##### 1.1.1 解决 N+1 查询
```go
// ❌ 当前问题代码
func GetDiariesWithAuthors(diaryIDs []int) ([]*Diary, error) {
    diaries := []*Diary{}
    for _, id := range diaryIDs {
        diary := GetDiary(id)
        diary.Author = GetUser(diary.AuthorID)  // N+1 查询
        diaries = append(diaries, diary)
    }
    return diaries, nil
}

// ✅ 优化后：使用 JOIN 查询
func GetDiariesWithAuthors(diaryIDs []int) ([]*Diary, error) {
    query := `
        SELECT
            d.*,
            u.id as author_id,
            u.name as author_name,
            u.email as author_email
        FROM diaries d
        LEFT JOIN users u ON d.author_id = u.id
        WHERE d.id IN (?)
    `
    query, args, _ := sqlx.In(query, diaryIDs)
    query = db.Rebind(query)

    rows, err := db.Query(query, args...)
    // ... 处理结果
}

// ✅ 或使用预加载模式
type DiaryRepository struct {
    db *sqlx.DB
}

func (r *DiaryRepository) FindWithPreload(opts ...PreloadOption) ([]*Diary, error) {
    preloader := NewPreloader(r.db)

    for _, opt := range opts {
        opt(preloader)
    }

    return preloader.Load()
}

// 使用示例
diaries, err := repo.FindWithPreload(
    WithAuthor(),
    WithComments(),
    WithLimit(20),
)
```

##### 1.1.2 添加缺失索引
```sql
-- 用户表索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- 日记表索引
CREATE INDEX idx_diaries_author_id ON diaries(author_id);
CREATE INDEX idx_diaries_created_at ON diaries(created_at);
CREATE INDEX idx_diaries_status ON diaries(status);

-- 复合索引
CREATE INDEX idx_diaries_author_status_created ON diaries(author_id, status, created_at);

-- 评论表索引
CREATE INDEX idx_comments_entry_id ON comments(entry_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);

-- 全文搜索索引
ALTER TABLE entries ADD FULLTEXT(title, content);
```

##### 1.1.3 查询性能监控
```go
// 实现查询性能监控中间件
type QueryLogger struct {
    threshold time.Duration
    logger    *slog.Logger
}

func (ql *QueryLogger) Wrap(db *sql.DB) *sql.DB {
    return sqlhooks.Wrap(db, &hooks{
        logger:    ql.logger,
        threshold: ql.threshold,
    })
}

type hooks struct {
    logger    *slog.Logger
    threshold time.Duration
}

func (h *hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
    return context.WithValue(ctx, "query_start", time.Now()), nil
}

func (h *hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
    start := ctx.Value("query_start").(time.Time)
    duration := time.Since(start)

    if duration > h.threshold {
        h.logger.Warn("slow query detected",
            slog.String("query", query),
            slog.Duration("duration", duration),
            slog.Any("args", args),
        )
    }

    // 记录到 Prometheus
    queryDuration.WithLabelValues(extractQueryType(query)).Observe(duration.Seconds())

    return ctx, nil
}
```

### 1.2 连接池优化

```go
// 优化数据库连接池配置
func NewDB(config *config.Database) (*sql.DB, error) {
    db, err := sql.Open("mysql", config.DSN)
    if err != nil {
        return nil, err
    }

    // 连接池配置
    db.SetMaxOpenConns(25)                  // 最大打开连接数
    db.SetMaxIdleConns(10)                  // 最大空闲连接数
    db.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生命周期
    db.SetConnMaxIdleTime(2 * time.Minute)  // 空闲连接最大生命周期

    // 验证连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, err
    }

    return db, nil
}
```

## 2. 缓存优化

### 2.1 多级缓存架构

```go
// 实现多级缓存
type CacheManager struct {
    l1Cache   *ristretto.Cache  // L1: 进程内缓存
    l2Cache   *redis.Client     // L2: Redis 缓存
    logger    *slog.Logger
}

func NewCacheManager(redis *redis.Client) (*CacheManager, error) {
    // 配置 L1 缓存
    l1Cache, err := ristretto.NewCache(&ristretto.Config{
        NumCounters: 10000,    // key 数量的 10 倍
        MaxCost:     100 << 20, // 100MB
        BufferItems: 64,
    })
    if err != nil {
        return nil, err
    }

    return &CacheManager{
        l1Cache: l1Cache,
        l2Cache: redis,
        logger:  slog.Default(),
    }, nil
}

// 获取缓存
func (cm *CacheManager) Get(ctx context.Context, key string) (interface{}, error) {
    // 1. 尝试从 L1 获取
    if value, found := cm.l1Cache.Get(key); found {
        cacheHitRate.WithLabelValues("l1").Inc()
        return value, nil
    }

    // 2. 尝试从 L2 获取
    cmd := cm.l2Cache.Get(ctx, key)
    if err := cmd.Err(); err == nil {
        value := cmd.Val()

        // 回填 L1
        cm.l1Cache.Set(key, value, 1)

        cacheHitRate.WithLabelValues("l2").Inc()
        return value, nil
    }

    cacheMissRate.Inc()
    return nil, ErrCacheMiss
}

// 设置缓存（写穿策略）
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    // 1. 写入 L1
    cm.l1Cache.SetWithTTL(key, value, 1, ttl)

    // 2. 异步写入 L2
    go func() {
        if err := cm.l2Cache.Set(context.Background(), key, value, ttl).Err(); err != nil {
            cm.logger.Error("failed to set L2 cache",
                slog.String("key", key),
                slog.Any("error", err),
            )
        }
    }()

    return nil
}
```

### 2.2 缓存预热

```go
// 应用启动时预热缓存
func PrewarmCache(ctx context.Context, cache *CacheManager, repo Repository) error {
    g, ctx := errgroup.WithContext(ctx)

    // 预热热门用户
    g.Go(func() error {
        users, err := repo.GetPopularUsers(ctx, 100)
        if err != nil {
            return err
        }

        for _, user := range users {
            key := fmt.Sprintf("user:%d", user.ID)
            cache.Set(ctx, key, user, 1*time.Hour)
        }
        return nil
    })

    // 预热热门日记
    g.Go(func() error {
        diaries, err := repo.GetPopularDiaries(ctx, 50)
        if err != nil {
            return err
        }

        for _, diary := range diaries {
            key := fmt.Sprintf("diary:%d", diary.ID)
            cache.Set(ctx, key, diary, 30*time.Minute)
        }
        return nil
    })

    // 预热配置数据
    g.Go(func() error {
        configs, err := repo.GetAllConfigs(ctx)
        if err != nil {
            return err
        }

        cache.Set(ctx, "configs:all", configs, 24*time.Hour)
        return nil
    })

    return g.Wait()
}
```

### 2.3 缓存更新策略

```go
// 缓存更新器
type CacheUpdater struct {
    cache *CacheManager
    repo  Repository
}

// 延迟双删策略
func (cu *CacheUpdater) UpdateUser(ctx context.Context, user *User) error {
    key := fmt.Sprintf("user:%d", user.ID)

    // 1. 删除缓存
    cu.cache.Delete(ctx, key)

    // 2. 更新数据库
    if err := cu.repo.UpdateUser(ctx, user); err != nil {
        return err
    }

    // 3. 延迟再次删除（防止并发问题）
    time.AfterFunc(500*time.Millisecond, func() {
        cu.cache.Delete(context.Background(), key)
    })

    return nil
}

// 主动更新策略
func (cu *CacheUpdater) RefreshUser(ctx context.Context, userID int) error {
    user, err := cu.repo.GetUser(ctx, userID)
    if err != nil {
        return err
    }

    key := fmt.Sprintf("user:%d", userID)
    return cu.cache.Set(ctx, key, user, 1*time.Hour)
}
```

## 3. API 性能优化

### 3.1 响应压缩

```go
// Gzip 压缩中间件
func GzipMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
            c.Next()
            return
        }

        gz := gzip.NewWriter(c.Writer)
        defer gz.Close()

        c.Header("Content-Encoding", "gzip")
        c.Writer = &gzipWriter{
            ResponseWriter: c.Writer,
            Writer:         gz,
        }

        c.Next()
    }
}

// Brotli 压缩支持
func BrotliMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if strings.Contains(c.GetHeader("Accept-Encoding"), "br") {
            br := brotli.NewWriter(c.Writer)
            defer br.Close()

            c.Header("Content-Encoding", "br")
            c.Writer = &brotliWriter{
                ResponseWriter: c.Writer,
                Writer:         br,
            }
        }
        c.Next()
    }
}
```

### 3.2 请求合并

```go
// 批量请求处理器
type BatchHandler struct {
    handlers map[string]HandlerFunc
    timeout  time.Duration
}

func (bh *BatchHandler) Handle(c *gin.Context) {
    var requests []BatchRequest
    if err := c.ShouldBindJSON(&requests); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    responses := make([]BatchResponse, len(requests))
    var wg sync.WaitGroup

    for i, req := range requests {
        wg.Add(1)
        go func(idx int, r BatchRequest) {
            defer wg.Done()

            ctx, cancel := context.WithTimeout(c.Request.Context(), bh.timeout)
            defer cancel()

            handler, ok := bh.handlers[r.Method]
            if !ok {
                responses[idx] = BatchResponse{
                    Error: "unknown method",
                }
                return
            }

            result, err := handler(ctx, r.Params)
            responses[idx] = BatchResponse{
                ID:     r.ID,
                Result: result,
                Error:  err,
            }
        }(i, req)
    }

    wg.Wait()
    c.JSON(200, responses)
}
```

### 3.3 GraphQL 优化（如果使用）

```go
// DataLoader 防止 N+1 查询
type UserLoader struct {
    *dataloader.Loader
}

func NewUserLoader(repo UserRepository) *UserLoader {
    return &UserLoader{
        Loader: dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
            userIDs := make([]int, len(keys))
            for i, key := range keys {
                userIDs[i] = key.Raw().(int)
            }

            users, err := repo.GetUsersByIDs(ctx, userIDs)
            if err != nil {
                return handleError(err, len(keys))
            }

            userMap := make(map[int]*User)
            for _, user := range users {
                userMap[user.ID] = user
            }

            results := make([]*dataloader.Result, len(keys))
            for i, id := range userIDs {
                user, ok := userMap[id]
                if ok {
                    results[i] = &dataloader.Result{Data: user}
                } else {
                    results[i] = &dataloader.Result{Error: errors.New("user not found")}
                }
            }

            return results
        }),
    }
}
```

## 4. 前端性能优化

### 4.1 资源加载优化

```javascript
// 实现资源预加载
class ResourcePreloader {
    constructor() {
        this.preloadedResources = new Set();
    }

    preloadImage(url) {
        if (this.preloadedResources.has(url)) return;

        const link = document.createElement('link');
        link.rel = 'preload';
        link.as = 'image';
        link.href = url;
        document.head.appendChild(link);

        this.preloadedResources.add(url);
    }

    preloadScript(url) {
        if (this.preloadedResources.has(url)) return;

        const link = document.createElement('link');
        link.rel = 'modulepreload';
        link.href = url;
        document.head.appendChild(link);

        this.preloadedResources.add(url);
    }

    // 预加载关键资源
    preloadCriticalResources() {
        // 预加载字体
        this.preloadFont('/fonts/main.woff2');

        // 预加载关键图片
        this.preloadImage('/images/logo.webp');

        // 预连接到 API 服务器
        this.preconnect('https://api.idiary.ai');
    }

    preconnect(url) {
        const link = document.createElement('link');
        link.rel = 'preconnect';
        link.href = url;
        document.head.appendChild(link);
    }
}
```

### 4.2 图片优化

```javascript
// 实现图片懒加载
class LazyImageLoader {
    constructor() {
        this.imageObserver = null;
        this.init();
    }

    init() {
        if ('IntersectionObserver' in window) {
            this.imageObserver = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        this.loadImage(entry.target);
                        this.imageObserver.unobserve(entry.target);
                    }
                });
            }, {
                rootMargin: '50px 0px',  // 提前 50px 加载
                threshold: 0.01
            });

            // 观察所有懒加载图片
            document.querySelectorAll('img[data-src]').forEach(img => {
                this.imageObserver.observe(img);
            });
        } else {
            // 降级方案
            this.loadAllImages();
        }
    }

    loadImage(img) {
        const src = img.dataset.src;
        const srcset = img.dataset.srcset;

        // 创建新图片对象预加载
        const tempImg = new Image();

        tempImg.onload = () => {
            img.src = src;
            if (srcset) img.srcset = srcset;
            img.classList.add('loaded');
        };

        tempImg.src = src;
    }

    // WebP 支持检测
    supportsWebP() {
        const canvas = document.createElement('canvas');
        canvas.width = 1;
        canvas.height = 1;
        return canvas.toDataURL('image/webp').indexOf('image/webp') === 0;
    }
}
```

### 4.3 JavaScript 优化

```javascript
// 实现代码分割和按需加载
class ModuleLoader {
    constructor() {
        this.loadedModules = new Map();
    }

    async load(moduleName) {
        // 检查缓存
        if (this.loadedModules.has(moduleName)) {
            return this.loadedModules.get(moduleName);
        }

        try {
            const module = await import(`./modules/${moduleName}.js`);
            this.loadedModules.set(moduleName, module);
            return module;
        } catch (error) {
            console.error(`Failed to load module: ${moduleName}`, error);
            throw error;
        }
    }

    // 路由级别的代码分割
    async loadRoute(routeName) {
        switch (routeName) {
            case 'diary':
                return this.load('diary-module');
            case 'user':
                return this.load('user-module');
            case 'admin':
                // 管理员模块按需加载
                return this.load('admin-module');
            default:
                throw new Error(`Unknown route: ${routeName}`);
        }
    }
}

// 防抖和节流优化
const utils = {
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },

    throttle(func, limit) {
        let inThrottle;
        return function(...args) {
            if (!inThrottle) {
                func.apply(this, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }
};

// 使用示例
const searchInput = document.getElementById('search');
searchInput.addEventListener('input', utils.debounce((e) => {
    performSearch(e.target.value);
}, 300));

window.addEventListener('scroll', utils.throttle(() => {
    updateScrollPosition();
}, 100));
```

## 5. 监控和分析

### 5.1 性能指标收集

```go
// Prometheus 指标定义
var (
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request latencies in seconds.",
            Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
        },
        []string{"method", "endpoint", "status"},
    )

    dbQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "db_query_duration_seconds",
            Help: "Database query latencies in seconds.",
            Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1.0},
        },
        []string{"query_type", "table"},
    )

    cacheHitRate = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_hit_total",
            Help: "Number of cache hits.",
        },
        []string{"cache_type"},
    )

    cacheMissRate = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "cache_miss_total",
            Help: "Number of cache misses.",
        },
    )
)

// 性能监控中间件
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        endpoint := c.FullPath()
        method := c.Request.Method

        httpRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
    }
}
```

### 5.2 性能分析工具

```go
// CPU 性能分析
func EnableProfiling(r *gin.Engine) {
    r.GET("/debug/pprof/", gin.WrapF(pprof.Index))
    r.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
    r.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
    r.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
    r.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))
    r.GET("/debug/pprof/heap", gin.WrapH(pprof.Handler("heap")))
    r.GET("/debug/pprof/goroutine", gin.WrapH(pprof.Handler("goroutine")))
    r.GET("/debug/pprof/allocs", gin.WrapH(pprof.Handler("allocs")))
}

// 自定义性能报告
func PerformanceReport() gin.HandlerFunc {
    return func(c *gin.Context) {
        report := map[string]interface{}{
            "goroutines": runtime.NumGoroutine(),
            "memory": getMemoryStats(),
            "database": getDatabaseStats(),
            "cache": getCacheStats(),
        }

        c.JSON(200, report)
    }
}

func getMemoryStats() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return map[string]interface{}{
        "alloc_mb":      m.Alloc / 1024 / 1024,
        "total_alloc_mb": m.TotalAlloc / 1024 / 1024,
        "sys_mb":        m.Sys / 1024 / 1024,
        "num_gc":        m.NumGC,
    }
}
```

## 6. 实施计划

### 第一阶段（第 1 周）
1. **数据库索引优化**
   - 添加缺失索引
   - 分析慢查询日志
   - 优化查询语句

2. **基础缓存实现**
   - 部署 Redis 缓存
   - 实现基本缓存逻辑
   - 缓存热点数据

### 第二阶段（第 2-3 周）
1. **查询优化**
   - 解决 N+1 查询问题
   - 实现查询结果缓存
   - 优化复杂查询

2. **前端优化**
   - 实现资源压缩
   - 图片懒加载
   - JavaScript 代码分割

### 第三阶段（第 4 周）
1. **高级缓存策略**
   - 实现多级缓存
   - 缓存预热机制
   - 缓存更新策略

2. **监控系统**
   - 部署 Prometheus
   - 配置 Grafana 仪表板
   - 设置性能告警

### 第四阶段（持续）
1. **持续优化**
   - 定期性能评估
   - A/B 测试
   - 用户反馈收集

2. **扩展优化**
   - CDN 部署
   - 服务拆分
   - 异步处理

## 7. 性能目标

### 短期目标（1 个月）
- API 响应时间 P95 < 200ms
- 页面加载时间 < 2s
- 数据库查询时间 P95 < 50ms

### 中期目标（3 个月）
- API 响应时间 P95 < 100ms
- 页面加载时间 < 1s
- 缓存命中率 > 80%

### 长期目标（6 个月）
- 支持 10x 当前用户量
- 99.9% 可用性
- 全球用户访问延迟 < 100ms

## 总结

通过实施这个性能优化方案，预期可以达到：

1. **响应时间减少 50%+**：通过查询优化和缓存策略
2. **数据库负载降低 70%+**：通过缓存和查询优化
3. **前端加载速度提升 60%+**：通过资源优化和懒加载
4. **系统可扩展性提升**：为未来增长做好准备

建议按照实施计划逐步推进，每个阶段都进行性能测试和监控，确保优化效果达到预期。