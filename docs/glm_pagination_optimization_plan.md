# 分页系统优化方案（Limit+1 模式）

## 一、现状分析

### 1.1 当前实现概述

项目已采用 **Limit+1** 分页模式，避免了 `COUNT(*)` 查询的性能问题。

**当前架构**：
```
Controller → Service → Repository
     ↓           ↓          ↓
  Paginator   HasMore    LIMIT pageSize+1
```

**已实现的功能**：
- ✅ Repository 层：`ListByUserID`、`ListByDiaryID` 等方法返回 `(items, hasMore, error)`
- ✅ Service 层：接口已适配 `hasMore` 返回值
- ✅ Controller 层：使用 `pkg.ListOptions` + `pkg.GeneratePagination`
- ✅ 前端模板：支持 HTMX 局部刷新，只显示"上一页/下一页"

### 1.2 当前问题

| 问题 | 位置 | 影响 |
|------|------|------|
| 遗留的 Count 方法 | `GetCountByUserID`、`GetCountByDiaryID` 等 | 代码冗余，可能被误用 |
| Service 层接口不一致 | 部分方法返回 `(items, count, error)` | 接口混乱 |
| 分页参数分散 | Controller 中手动解析 pageIndex | 代码重复 |
| 缺少分页元数据 | 无首页、末页判断 | 无法做更多优化 |
| 前端无法跳页 | 只有上/下一页 | 用户体验受限 |

---

## 二、优化目标

1. **统一分页接口**：所有分页方法统一返回 `(items, hasMore, error)`
2. **消除 Count 查询**：移除不再使用的 Count 方法
3. **封装分页逻辑**：创建通用的分页请求/响应结构
4. **增强分页能力**：支持游标分页、页码跳转等高级功能
5. **优化前端体验**：提供更丰富的分页交互

---

## 三、核心优化方案

### 3.1 统一分页结构体

**新增文件**：`pkg/paging/types.go`

```go
package paging

// PageRequest 分页请求参数
type PageRequest struct {
    PageIndex int    // 当前页码（从1开始）
    PageSize  int    // 每页大小
    Cursor    string // 游标（可选，用于游标分页）
}

// NewPageRequest 创建分页请求，使用默认值
func NewPageRequest(pageIndex, pageSize int) *PageRequest {
    if pageIndex < 1 {
        pageIndex = 1
    }
    if pageSize <= 0 {
        pageSize = 20 // 默认每页20条
    }
    return &PageRequest{
        PageIndex: pageIndex,
        PageSize:  pageSize,
    }
}

// Offset 计算 OFFSET 值
func (r *PageRequest) Offset() int {
    return (r.PageIndex - 1) * r.PageSize
}

// Limit 计算 LIMIT 值（Limit+1 模式）
func (r *PageRequest) Limit() int {
    return r.PageSize + 1
}

// PageResult 分页结果（泛型版本）
type PageResult[T any] struct {
    Items     []T    // 数据列表
    HasMore   bool   // 是否有更多数据
    PageIndex int    // 当前页码
    PageSize  int    // 每页大小
    HasPrev   bool   // 是否有上一页
}

// NewPageResult 创建分页结果
func NewPageResult[T any](items []T, pageSize, pageIndex int) *PageResult[T] {
    hasMore := len(items) > pageSize
    if hasMore {
        items = items[:pageSize] // 截取实际需要的数据
    }
    hasPrev := pageIndex > 1
    
    return &PageResult[T]{
        Items:     items,
        HasMore:   hasMore,
        PageIndex: pageIndex,
        PageSize:  pageSize,
        HasPrev:   hasPrev,
    }
}
```

### 3.2 优化 Paginator.go

**优化文件**：`pkg/Paginator.go`

```go
package pkg

import (
    "fmt"
)

// ListOptions 分页选项（limit+1 模式）
type ListOptions struct {
    PageIndex int    // 当前页码（从1开始）
    PageSize  int    // 每页显示的记录数
    HasMore   bool   // 是否还有更多数据（limit+1 模式）
    HasPrev   bool   // 是否有上一页
    URL       string // 分页链接，其中应该包含 %d，用于替换为页码
    PagePrev  string // 上一页文字
    PageNext  string // 下一页文字
    HxTarget  string // HTMX 目标元素（可选）
}

// PageLinks 分页链接结构
type PageLinks struct {
    URL       string // 链接地址
    Text      string // 链接文字
    PageIndex int    // 页码
    IsPrev    bool   // 是否为上一页
    IsNext    bool   // 是否为下一页
    HxTarget  string // HTMX 目标元素
}

// GeneratePagination 生成分页导航（HasMore 模式）
// 只显示上一页和下一页，适用于"加载更多"场景
func GeneratePagination(listOptions ListOptions) []PageLinks {
    pageSize := listOptions.PageSize
    pageIndex := listOptions.PageIndex

    // 参数校验
    if pageSize <= 0 {
        return []PageLinks{}
    }
    if pageIndex < 1 {
        pageIndex = 1
    }

    aryLinks := make([]PageLinks, 0)

    // 如果当前页码大于 1，则添加上一页链接
    if listOptions.HasPrev || pageIndex > 1 {
        aryLinks = append(aryLinks, PageLinks{
            URL:      fmt.Sprintf(listOptions.URL, pageIndex-1),
            Text:     listOptions.PagePrev,
            PageIndex: pageIndex - 1,
            IsPrev:   true,
            HxTarget: listOptions.HxTarget,
        })
    }

    // 如果还有更多数据，则添加下一页链接
    if listOptions.HasMore {
        aryLinks = append(aryLinks, PageLinks{
            URL:       fmt.Sprintf(listOptions.URL, pageIndex+1),
            Text:      listOptions.PageNext,
            PageIndex: pageIndex + 1,
            IsNext:    true,
            HxTarget:  listOptions.HxTarget,
        })
    }

    return aryLinks
}

// GeneratePageNumbers 生成分页数字链接（可选功能）
// 在 HasMore 模式下，提供有限的页码跳转能力
func GeneratePageNumbers(listOptions ListOptions, maxNumbers int) []PageLinks {
    if maxNumbers <= 0 {
        maxNumbers = 5
    }
    
    pageIndex := listOptions.PageIndex
    if pageIndex < 1 {
        pageIndex = 1
    }

    links := make([]PageLinks, 0)
    
    // 计算显示的页码范围
    start := pageIndex - maxNumbers/2
    if start < 1 {
        start = 1
    }
    
    end := start + maxNumbers - 1
    
    for i := start; i <= end; i++ {
        isCurrent := i == pageIndex
        // 已访问过的页面总是可以跳转
        // 未访问的页面只有在 HasMore 且 i <= pageIndex+1 时才显示
        if i <= pageIndex || (listOptions.HasMore && i == pageIndex+1) {
            links = append(links, PageLinks{
                URL:       fmt.Sprintf(listOptions.URL, i),
                Text:      fmt.Sprintf("%d", i),
                PageIndex: i,
                HxTarget:  listOptions.HxTarget,
            })
        }
        if isCurrent {
            // 标记当前页
            links[len(links)-1].Text = fmt.Sprintf("[%d]", i)
        }
    }
    
    return links
}
```

### 3.3 Repository 层优化

**统一接口**：`internal/domain/diary/repository.go`

```go
// ListByUserID 根据特定用户ID检索分页的日记本列表
// 参数:
//   - userID: 用户ID
//   - includePrivate: 是否包含私有日记本
//   - pageIndex: 页码（从1开始）
//   - pageSize: 每页大小
// 返回:
//   - []Diary: 日记本列表（实际返回 pageSize 条）
//   - hasMore: 是否还有更多数据
//   - error: 查询失败时返回错误
ListByUserID(userID int64, includePrivate bool, pageIndex, pageSize int) ([]Diary, bool, error)

// 移除：不再需要 Count 方法
// GetCountByUserID(userID int64, includePrivate bool) (int, error)  // 删除
```

**Repository 实现**：`internal/infra/repository/diary.go`

```go
// ListByUserID 根据特定用户ID检索分页的日记本列表
// 使用 Limit+1 模式判断是否有下一页
func (r *DiaryRepository) ListByUserID(userID int64, includePrivate bool, pageIndex, pageSize int) ([]diary.Diary, bool, error) {
    diaryArray := []diary.Diary{}
    
    // 参数校验
    if pageIndex < 1 {
        pageIndex = 1
    }
    if pageSize <= 0 {
        pageSize = 20
    }
    
    offset := (pageIndex - 1) * pageSize
    limit := pageSize + 1 // Limit+1 模式

    baseSQL := `SELECT d.*, 
               u1.handle, u1.display_name
               FROM br_diary d
               LEFT JOIN br_user u1 ON d.add_user_id = u1.user_id
               WHERE d.add_user_id = ?`

    var args []interface{}
    args = append(args, userID)

    if !includePrivate {
        baseSQL += " AND d.view_setting = ?"
        args = append(args, diary.ViewSettingPublic)
    }

    baseSQL += ` ORDER BY d.diary_id DESC LIMIT ? OFFSET ?`
    args = append(args, limit, offset)

    err := r.db.Select(&diaryArray, baseSQL, args...)
    if err == sql.ErrNoRows {
        return []diary.Diary{}, false, nil
    }
    if err != nil {
        return nil, false, err
    }

    // 判断是否还有更多数据
    hasMore := len(diaryArray) > pageSize
    if hasMore {
        diaryArray = diaryArray[:pageSize] // 截取实际需要的 pageSize 条
    }

    return diaryArray, hasMore, nil
}

// 删除以下方法（不再需要）：
// - GetCountByUserID
// - GetJoinedDiariesCountByUserID
// - GetActiveJoinedDiariesCountByUserID
```

### 3.4 Service 层优化

**优化文件**：`internal/domain/diary/service.go`

```go
// ListDiariesByUser 获取用户的分页日记本列表
// 参数:
//   - userID: 用户ID
//   - includePrivate: 是否包含私有日记本
//   - pageIndex: 页码（从1开始）
//   - pageSize: 每页大小
// 返回:
//   - []Diary: 日记本列表
//   - hasMore: 是否还有更多数据
//   - error: 查询失败时返回错误
func (s *Service) ListDiariesByUser(userID int64, includePrivate bool, pageIndex, pageSize int) ([]Diary, bool, error) {
    return s.repo.ListByUserID(userID, includePrivate, pageIndex, pageSize)
}

// 移除旧的返回 count 的方法
```

### 3.5 Controller 层优化

**创建分页辅助方法**：`internal/app/controller/helpers.go`

```go
package controller

import (
    "strconv"
    "strings"
    
    "github.com/biuren/hebe/pkg"
    "github.com/gin-gonic/gin"
)

// ParsePageParams 解析分页参数
func ParsePageParams(c *gin.Context, defaultPageSize int) (pageIndex, pageSize int) {
    pageIndex = 1
    pageSize = defaultPageSize
    
    // 从 URL 参数解析
    if pageStr := c.Query("page"); pageStr != "" {
        if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
            pageIndex = page
        }
    }
    
    // 从路径参数解析（如 /u/123/p/2）
    if pageStr := c.Param("pageindex"); pageStr != "" {
        if page, err := strconv.Atoi(strings.TrimSpace(pageStr)); err == nil && page > 0 {
            pageIndex = page
        }
    }
    
    return
}

// BuildPageLinks 构建分页链接
func BuildPageLinks(pageIndex, pageSize int, hasMore bool, urlTemplate, prevText, nextText, hxTarget string) []pkg.PageLinks {
    listOptions := pkg.ListOptions{
        PageIndex: pageIndex,
        PageSize:  pageSize,
        HasMore:   hasMore,
        HasPrev:   pageIndex > 1,
        URL:       urlTemplate,
        PagePrev:  prevText,
        PageNext:  nextText,
        HxTarget:  hxTarget,
    }
    return pkg.GeneratePagination(listOptions)
}
```

**Controller 使用示例**：`internal/app/controller/diary/view.go`

```go
func (ctrl *DiaryController) ShowUserProfile(c *gin.Context) {
    // 解析用户ID
    uid, err := validation.ParseID(c.Param("uid"))
    if err != nil {
        ctrl.renderService.Redirect(c, "/")
        return
    }

    // 解析分页参数（使用辅助方法）
    pageIndex, pageSize := controller.ParsePageParams(c, ctrl.config.App.PageSize)

    // 创建页面数据
    data := ctrl.renderService.CreatePageData(c, "user.profile.title")

    // 判断是否为页面所有者
    isOwner := data.IsAuthenticated() && data.GetAuthUser().UserId == uid

    // 获取用户的日记本列表
    diaries, hasMore, err := ctrl.diaryService.ListDiariesByUser(uid, isOwner, pageIndex, pageSize)
    if err != nil {
        ctrl.renderService.RenderErrorPage(c, errors.InternalError)
        return
    }

    // 构建分页链接（使用辅助方法）
    pageLinks := controller.BuildPageLinks(
        pageIndex, pageSize, hasMore,
        fmt.Sprintf("/u/%d/p/%%d", uid),
        "<<", ">>", "#diary-container",
    )

    // 设置页面数据
    data.Data["Diaries"] = diaries
    data.Data["PageLinks"] = pageLinks
    data.Data["PageIndex"] = pageIndex
    data.Data["PageSize"] = pageSize
    data.Data["HasMore"] = hasMore
    data.Data["HxTarget"] = "#diary-container"

    // HTMX 请求处理
    if c.GetHeader("HX-Request") == "true" {
        ctrl.renderService.HTML(c, "user/list", data)
        return
    }

    ctrl.renderService.HTML(c, "user/index", data)
}
```

---

## 四、高级优化：游标分页（可选）

对于大数据量场景，可以使用**游标分页**替代**偏移分页**，性能更优。

### 4.1 游标分页结构

```go
// pkg/paging/cursor.go

import (
    "encoding/base64"
    "encoding/json"
)

// CursorRequest 游标分页请求
type CursorRequest struct {
    Cursor string // 上次查询的游标
    Limit  int    // 每页大小
}

// CursorResult 游标分页结果
type CursorResult[T any] struct {
    Items       []T    // 数据列表
    NextCursor  string // 下一页游标（空表示无更多数据）
    HasMore     bool   // 是否有更多数据
}

// EncodeCursor 编码游标（基于最后一条记录的ID）
func EncodeCursor(lastID int64) string {
    data, _ := json.Marshal(map[string]int64{"id": lastID})
    return base64.URLEncoding.EncodeToString(data)
}

// DecodeCursor 解码游标
func DecodeCursor(cursor string) (int64, error) {
    data, err := base64.URLEncoding.DecodeString(cursor)
    if err != nil {
        return 0, err
    }
    var result map[string]int64
    if err := json.Unmarshal(data, &result); err != nil {
        return 0, err
    }
    return result["id"], nil
}
```

### 4.2 Repository 层实现

```go
// ListByUserIDWithCursor 使用游标分页获取日记本列表
func (r *DiaryRepository) ListByUserIDWithCursor(userID int64, cursor string, limit int) ([]diary.Diary, string, error) {
    diaryArray := []diary.Diary{}
    
    if limit <= 0 {
        limit = 20
    }
    
    baseSQL := `SELECT d.*, u1.handle, u1.display_name
               FROM br_diary d
               LEFT JOIN br_user u1 ON d.add_user_id = u1.user_id
               WHERE d.add_user_id = ? AND d.status = ?`
    
    var args []interface{}
    args = append(args, userID, diary.StatusActive)
    
    // 解码游标，获取起始ID
    if cursor != "" {
        lastID, err := paging.DecodeCursor(cursor)
        if err == nil && lastID > 0 {
            baseSQL += " AND d.diary_id < ?" // 降序排列，使用小于
            args = append(args, lastID)
        }
    }
    
    baseSQL += " ORDER BY d.diary_id DESC LIMIT ?"
    args = append(args, limit+1) // Limit+1
    
    err := r.db.Select(&diaryArray, baseSQL, args...)
    if err != nil {
        return nil, "", err
    }
    
    // 判断是否有更多数据
    var nextCursor string
    if len(diaryArray) > limit {
        nextCursor = paging.EncodeCursor(diaryArray[limit-1].DiaryID)
        diaryArray = diaryArray[:limit]
    }
    
    return diaryArray, nextCursor, nil
}
```

---

## 五、前端优化

### 5.1 分页模板优化

**优化文件**：`website/views/pagination.html`

```html
{{define "pagination"}}
{{/* Limit+1 模式：只显示 prev/next 导航 */}}
{{if .Data.PageLinks}}
<div class="pagination-container flex items-center justify-between">
    <!-- 页码信息 -->
    <div class="text-sm text-gray-500">
        {{if .Data.HasMore}}
        {{getMessage "pagination.has_more" .Context.Lang}}
        {{else}}
        {{getMessage "pagination.no_more" .Context.Lang}}
        {{end}}
    </div>
    
    <!-- 导航按钮 -->
    <div class="flex space-x-2">
        {{range .Data.PageLinks}}
        {{if $.Data.HxTarget}}
        <a class="pagination-btn border py-2 px-4 border-gray-300 dark:border-gray-600 
                  hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300
                  {{if .IsPrev}}rounded-l{{end}}{{if .IsNext}}rounded-r{{end}}"
           href="{{.URL}}" 
           hx-get="{{.URL}}" 
           hx-target="{{$.Data.HxTarget}}" 
           hx-swap="outerHTML"
           hx-indicator="#loading-spinner">
            {{if .IsPrev}}<i class="mdi mdi-chevron-left"></i>{{end}}
            {{.Text}}
            {{if .IsNext}}<i class="mdi mdi-chevron-right"></i>{{end}}
        </a>
        {{else}}
        <a class="pagination-btn border py-2 px-4 border-gray-300 dark:border-gray-600 
                  hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300"
           href="{{.URL}}">
            {{if .IsPrev}}<i class="mdi mdi-chevron-left"></i>{{end}}
            {{.Text}}
            {{if .IsNext}}<i class="mdi mdi-chevron-right"></i>{{end}}
        </a>
        {{end}}
        {{end}}
    </div>
</div>

<!-- 加载指示器 -->
<div id="loading-spinner" class="htmx-indicator flex justify-center py-4">
    <i class="mdi mdi-loading mdi-spin text-2xl text-gray-400"></i>
</div>
{{end}}
{{end}}
```

### 5.2 无限滚动支持

```html
{{define "pagination-infinite"}}
{{/* 无限滚动模式 */}}
{{if .Data.HasMore}}
<div class="load-more-container py-4" 
     hx-get="{{printf .Data.PageURL (add .Data.PageIndex 1)}}" 
     hx-target="#list-container" 
     hx-swap="beforeend"
     hx-trigger="revealed"
     hx-indicator="#load-more-spinner">
    <div id="load-more-spinner" class="htmx-indicator flex justify-center">
        <i class="mdi mdi-loading mdi-spin text-2xl text-gray-400"></i>
    </div>
</div>
{{end}}
{{end}}
```

### 5.3 CSS 样式

```css
/* 分页按钮样式 */
.pagination-btn {
    transition: all 0.2s ease;
}

.pagination-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.pagination-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

/* HTMX 加载指示器 */
.htmx-indicator {
    opacity: 0;
    transition: opacity 200ms ease-in;
}

.htmx-request .htmx-indicator,
.htmx-request.htmx-indicator {
    opacity: 1;
}
```

---

## 六、需要删除的遗留代码

### 6.1 Repository 层

| 文件 | 方法 | 说明 |
|------|------|------|
| `diary.go` | `GetCountByUserID` | 不再使用 |
| `diary_author.go` | `GetJoinedDiariesCountByUserID` | 不再使用 |
| `diary_author.go` | `GetActiveJoinedDiariesCountByUserID` | 仅在限制检查时使用，可优化 |
| `entry.go` | `GetCountByUserID` | 不再使用 |
| `entry.go` | `GetCountByDiaryID` | 不再使用 |

### 6.2 Service 层

| 文件 | 方法 | 说明 |
|------|------|------|
| `diary/service.go` | `ListDiariesByUser` 返回 `(items, count, error)` | 改为返回 `(items, hasMore, error)` |

### 6.3 接口定义

| 文件 | 接口 | 说明 |
|------|------|------|
| `diary/repository.go` | `GetCountByUserID` | 删除 |
| `diary/repository.go` | `GetJoinedDiariesCountByUserID` | 删除 |
| `entry/repository.go` | `GetCountByUserID` | 删除 |
| `entry/repository.go` | `GetCountByDiaryID` | 删除 |

---

## 七、实施计划

### 阶段一：基础优化（1-2天）

1. ✅ 创建统一的分页结构体 `pkg/paging/types.go`
2. ✅ 更新 `pkg/Paginator.go`，添加 `HasPrev` 支持
3. ✅ 创建 Controller 分页辅助方法
4. ✅ 更新前端模板

### 阶段二：接口统一（2-3天）

1. ⬜ 更新 Repository 接口定义
2. ⬜ 更新 Service 层方法签名
3. ⬜ 删除遗留的 Count 方法
4. ⬜ 更新所有 Controller 调用

### 阶段三：测试与验证（1-2天）

1. ⬜ 更新单元测试
2. ⬜ 更新集成测试
3. ⬜ 手动测试所有分页功能
4. ⬜ 性能测试对比

### 阶段四：高级功能（可选）

1. ⬜ 实现游标分页
2. ⬜ 实现无限滚动
3. ⬜ 添加页码跳转

---

## 八、性能对比

### 优化前（Count 分页）

```sql
-- 查询1：获取总数
SELECT COUNT(*) FROM br_diary WHERE add_user_id = ? AND status = ?

-- 查询2：获取数据
SELECT * FROM br_diary WHERE add_user_id = ? AND status = ? 
ORDER BY diary_id DESC LIMIT 20 OFFSET 0
```

**问题**：COUNT(*) 在大数据量时性能较差。

### 优化后（Limit+1 分页）

```sql
-- 仅需一次查询
SELECT * FROM br_diary WHERE add_user_id = ? AND status = ? 
ORDER BY diary_id DESC LIMIT 21 OFFSET 0
```

**优势**：
- 查询次数减少 50%
- 无需 COUNT(*) 扫描
- 大数据量时性能稳定

---

## 九、注意事项

### 9.1 分页一致性

在数据频繁变化的场景下，Limit+1 分页可能出现：
- **数据重复**：新增数据导致同一数据出现在不同页
- **数据遗漏**：删除数据导致部分数据被跳过

**解决方案**：
- 对于实时性要求高的场景，使用游标分页
- 对于一般场景，Limit+1 分页已足够

### 9.2 首页特殊处理

首页（pageIndex=1）不需要显示"上一页"按钮，通过 `HasPrev` 判断。

### 9.3 空数据处理

当返回空列表时，`HasMore` 应为 `false`，模板应显示"暂无数据"提示。

---

## 十、总结

本优化方案通过以下措施提升分页性能和代码质量：

| 优化项 | 效果 |
|--------|------|
| 移除 COUNT 查询 | 减少 50% 数据库查询 |
| 统一分页接口 | 提高代码一致性和可维护性 |
| 封装分页逻辑 | 减少重复代码 |
| 游标分页（可选） | 大数据量场景性能更优 |
| 前端优化 | 提升用户体验 |

当前项目已实现了大部分 Limit+1 分页逻辑，优化工作主要是**清理遗留代码**和**统一接口规范**。

---

*方案生成时间: 2026-03-04*
