package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	contracts "github.com/biuren/hebe/internal/contracts"
)

// TestAuthUser 测试 AuthUser 结构体
func TestAuthUser(t *testing.T) {
	t.Run("创建AuthUser实例", func(t *testing.T) {
		user := &contracts.AuthUser{
			LoginName:    "testuser",
			AccessUrl:    "https://example.com",
			UserId:       12345,
			Handle:       "testuser",
			DisplayName:  "Test User",
			Email:        "test@example.com",
			UserGroup:    1,
			FaceUrl:      "https://example.com/avatar.jpg",
			LoginUTCTime: 1640995200000,
			TimeZone:     "UTC",
			Lang:         "en",
			Status:       1,
		}

		assert.Equal(t, "testuser", user.LoginName)
		assert.Equal(t, "https://example.com", user.AccessUrl)
		assert.Equal(t, int64(12345), user.UserId)
		assert.Equal(t, "testuser", user.Handle)
		assert.Equal(t, "Test User", user.DisplayName)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, int8(1), user.UserGroup)
		assert.Equal(t, "https://example.com/avatar.jpg", user.FaceUrl)
		assert.Equal(t, int64(1640995200000), user.LoginUTCTime)
		assert.Equal(t, "UTC", user.TimeZone)
		assert.Equal(t, "en", user.Lang)
		assert.Equal(t, int8(1), user.Status)
	})
}

// TestRenderContext 测试 RenderContext 结构体
func TestRenderContext(t *testing.T) {
	t.Run("创建RenderContext实例", func(t *testing.T) {
		user := &contracts.AuthUser{
			LoginName: "testuser",
			UserId:    12345,
			Email:     "test@example.com",
		}

		ctx := &contracts.RenderContext{
			IsLogin:  true,
			AuthUser: user,
			Lang:     "en",
			Theme:    "default",
		}

		assert.True(t, ctx.IsLogin)
		assert.NotNil(t, ctx.AuthUser)
		assert.Equal(t, "testuser", ctx.AuthUser.LoginName)
		assert.Equal(t, "en", ctx.Lang)
		assert.Equal(t, "default", ctx.Theme)
	})

	t.Run("未登录状态的RenderContext", func(t *testing.T) {
		ctx := &contracts.RenderContext{
			IsLogin:  false,
			AuthUser: nil,
			Lang:     "zh-CN",
		}

		assert.False(t, ctx.IsLogin)
		assert.Nil(t, ctx.AuthUser)
		assert.Equal(t, "zh-CN", ctx.Lang)
	})
}

// TestPageData 测试 PageData 结构体和方法
func TestPageData(t *testing.T) {
	t.Run("创建PageData实例", func(t *testing.T) {
		user := &contracts.AuthUser{
			LoginName: "testuser",
			UserId:    12345,
			Email:     "test@example.com",
		}

		ctx := &contracts.RenderContext{
			IsLogin:  true,
			AuthUser: user,
			Lang:     "en",
		}

		data := map[string]interface{}{
			"entries": []string{"entry1", "entry2"},
			"count":   2,
		}

		pageData := &contracts.PageData[map[string]interface{}]{
			Title:       "Test Page",
			Description: "This is a test page",
			Keywords:    "test, page, example",
			Context:     ctx,
			Data:        data,
		}

		assert.Equal(t, "Test Page", pageData.Title)
		assert.Equal(t, "This is a test page", pageData.Description)
		assert.Equal(t, "test, page, example", pageData.Keywords)
		assert.NotNil(t, pageData.Context)
		assert.NotNil(t, pageData.Data)
		assert.Equal(t, 2, pageData.Data["count"])
	})

	t.Run("测试IsAuthenticated方法 - 已认证", func(t *testing.T) {
		user := &contracts.AuthUser{
			LoginName: "testuser",
			UserId:    12345,
		}

		ctx := &contracts.RenderContext{
			IsLogin:  true,
			AuthUser: user,
			Lang:     "en",
		}

		pageData := &contracts.PageData[map[string]interface{}]{
			Title:   "Test Page",
			Context: ctx,
			Data:    make(map[string]interface{}),
		}

		assert.True(t, pageData.IsAuthenticated())
	})

	t.Run("测试IsAuthenticated方法 - 未认证", func(t *testing.T) {
		ctx := &contracts.RenderContext{
			IsLogin:  false,
			AuthUser: nil,
			Lang:     "en",
		}

		pageData := &contracts.PageData[map[string]interface{}]{
			Title:   "Test Page",
			Context: ctx,
			Data:    make(map[string]interface{}),
		}

		assert.False(t, pageData.IsAuthenticated())
	})

	t.Run("测试IsAuthenticated方法 - Context为nil", func(t *testing.T) {
		pageData := &contracts.PageData[map[string]interface{}]{
			Title:   "Test Page",
			Context: nil,
			Data:    make(map[string]interface{}),
		}

		assert.False(t, pageData.IsAuthenticated())
	})

	t.Run("测试GetAuthUser方法 - 有用户", func(t *testing.T) {
		user := &contracts.AuthUser{
			LoginName: "testuser",
			UserId:    12345,
			Email:     "test@example.com",
		}

		ctx := &contracts.RenderContext{
			IsLogin:  true,
			AuthUser: user,
			Lang:     "en",
		}

		pageData := &contracts.PageData[map[string]interface{}]{
			Title:   "Test Page",
			Context: ctx,
			Data:    make(map[string]interface{}),
		}

		authUser := pageData.GetAuthUser()
		require.NotNil(t, authUser)
		assert.Equal(t, "testuser", authUser.LoginName)
		assert.Equal(t, int64(12345), authUser.UserId)
		assert.Equal(t, "test@example.com", authUser.Email)
	})

	t.Run("测试GetAuthUser方法 - 无用户", func(t *testing.T) {
		ctx := &contracts.RenderContext{
			IsLogin:  false,
			AuthUser: nil,
			Lang:     "en",
		}

		pageData := &contracts.PageData[map[string]interface{}]{
			Title:   "Test Page",
			Context: ctx,
			Data:    make(map[string]interface{}),
		}

		authUser := pageData.GetAuthUser()
		assert.Nil(t, authUser)
	})

	t.Run("测试GetAuthUser方法 - Context为nil", func(t *testing.T) {
		pageData := &contracts.PageData[map[string]interface{}]{
			Title:   "Test Page",
			Context: nil,
			Data:    make(map[string]interface{}),
		}

		authUser := pageData.GetAuthUser()
		assert.Nil(t, authUser)
	})
}

// TestPageDataWithDifferentTypes 测试 PageData 泛型功能
func TestPageDataWithDifferentTypes(t *testing.T) {
	t.Run("PageData with string data", func(t *testing.T) {
		ctx := &contracts.RenderContext{
			IsLogin: false,
			Lang:    "en",
		}

		pageData := &contracts.PageData[string]{
			Title:   "String Page",
			Context: ctx,
			Data:    "Hello, World!",
		}

		assert.Equal(t, "String Page", pageData.Title)
		assert.Equal(t, "Hello, World!", pageData.Data)
		assert.False(t, pageData.IsAuthenticated())
	})

	t.Run("PageData with struct data", func(t *testing.T) {
		type TestData struct {
			Name  string
			Count int
		}

		ctx := &contracts.RenderContext{
			IsLogin: false,
			Lang:    "en",
		}

		testData := TestData{
			Name:  "Test",
			Count: 42,
		}

		pageData := &contracts.PageData[TestData]{
			Title:   "Struct Page",
			Context: ctx,
			Data:    testData,
		}

		assert.Equal(t, "Struct Page", pageData.Title)
		assert.Equal(t, "Test", pageData.Data.Name)
		assert.Equal(t, 42, pageData.Data.Count)
	})

	t.Run("PageData with slice data", func(t *testing.T) {
		ctx := &contracts.RenderContext{
			IsLogin: false,
			Lang:    "en",
		}

		sliceData := []int{1, 2, 3, 4, 5}

		pageData := &contracts.PageData[[]int]{
			Title:   "Slice Page",
			Context: ctx,
			Data:    sliceData,
		}

		assert.Equal(t, "Slice Page", pageData.Title)
		assert.Len(t, pageData.Data, 5)
		assert.Equal(t, 1, pageData.Data[0])
		assert.Equal(t, 5, pageData.Data[4])
	})
}