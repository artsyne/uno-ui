package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	post "github.com/biuren/hebe/internal/app/controller/post"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/like"
	postDomain "github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// TestNewPostController 测试控制器创建
func TestNewPostController(t *testing.T) {
	cfg := &config.Config{}
	mockPostSvc := &postDomain.Service{}
	mockDiarySvc := &diary.Service{}
	mockAuth := &user.AuthHelper{}
	mockViewSvc := &view.ViewService{}
	mockUserSvc := &user.Service{}
	var mockLikeSvc like.Service

	// 创建控制器
	controller := post.NewPostController(
		cfg,
		mockPostSvc,
		mockDiarySvc,
		mockAuth,
		mockViewSvc,
		mockUserSvc,
		mockLikeSvc,
	)

	// 验证控制器不为空
	assert.NotNil(t, controller)
}

// TestDeletePost 测试删除日记请求解析
func TestDeletePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		uriID       string
		requestBody string
		contentType string
		expectError bool
		errorMsg    string
		expectedID  int64
	}{
		{
			name:        "有效的删除请求",
			uriID:       "123",
			requestBody: `{"diary_id": 456}`,
			contentType: "application/json",
			expectError: false,
			errorMsg:    "",
			expectedID:  123,
		},
		{
			name:        "无效的post_id",
			uriID:       "abc",
			requestBody: `{"diary_id": 456}`,
			contentType: "application/json",
			expectError: true,
			errorMsg:    "strconv.ParseInt", // ShouldBindUri returns error for type mismatch
			expectedID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建HTTP请求
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/post/"+tt.uriID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", tt.contentType)

			// 创建Gin上下文
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.uriID}}

			// 测试请求解析逻辑
			var deleteReq post.DeletePostRequest
			err := c.ShouldBindUri(&deleteReq)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_ = c.ShouldBind(&deleteReq)
				assert.Equal(t, tt.expectedID, deleteReq.PostID)
			}
		})
	}
}

// TestDeletePostFormData 测试表单数据删除日记
func TestDeletePostFormData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	// DELETE request usually takes query params, not body
	req, _ := http.NewRequest("DELETE", "/post/123?diary_id=456", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	var deleteReq post.DeletePostRequest
	err := c.ShouldBindUri(&deleteReq)
	assert.NoError(t, err)

	err = c.ShouldBind(&deleteReq) // ShouldBind binds query for DELETE
	assert.NoError(t, err)

	assert.Equal(t, int64(123), deleteReq.PostID)
	assert.Equal(t, int64(456), deleteReq.DiaryID)
}

// TestDeletePostValidation 测试删除请求验证
func TestDeletePostValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		postID      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效的post_id",
			postID:      "123",
			expectError: false,
			errorMsg:    "",
		},
		{
			name:        "非数字post_id",
			postID:      "abc",
			expectError: true,
			errorMsg:    "strconv.ParseInt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			// DELETE request with URI param

			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("DELETE", "/post/"+tt.postID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.postID}}

			var deleteReq post.DeletePostRequest
			err := c.ShouldBindUri(&deleteReq)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(123), deleteReq.PostID)
			}
		})
	}
}

// TestGetPostsStatsRequest 测试获取文章统计信息请求解析
func TestGetPostsStatsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		diaryID     string
		postIDs     string
		expectError bool
		errorMsg    string
		expectedLen int
	}{
		{
			name:        "有效的统计请求 - 单个ID",
			diaryID:     "123",
			postIDs:     "456",
			expectError: false,
			errorMsg:    "",
			expectedLen: 1,
		},
		{
			name:        "有效的统计请求 - 多个ID",
			diaryID:     "123",
			postIDs:     "456,789,101112",
			expectError: false,
			errorMsg:    "",
			expectedLen: 3,
		},
		{
			name:        "有效的统计请求 - 带空格的ID列表",
			diaryID:     "123",
			postIDs:     "456, 789 , 101112",
			expectError: false,
			errorMsg:    "",
			expectedLen: 3,
		},
		{
			name:        "无效的diaryID",
			diaryID:     "0",
			postIDs:     "456",
			expectError: true,
			errorMsg:    "required",
			expectedLen: 0,
		},
		{
			name:        "空的postIDs",
			diaryID:     "123",
			postIDs:     "",
			expectError: true,
			errorMsg:    "required",
			expectedLen: 0,
		},
		{
			name:        "无效的postID格式",
			diaryID:     "123",
			postIDs:     "456,abc,789",
			expectError: true,
			errorMsg:    "invalid post id",
			expectedLen: 0,
		},
		{
			name:        "包含负数的postID",
			diaryID:     "123",
			postIDs:     "456,-1,789",
			expectError: false, // ParseInt可以解析负数
			errorMsg:    "",
			expectedLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路径参数
			c.Params = gin.Params{
				{Key: "diaryid", Value: tt.diaryID},
			}

			// 设置查询参数
			req, _ := http.NewRequest("GET", "/d/"+tt.diaryID+"/posts/stats?postids="+tt.postIDs, nil)
			c.Request = req

			// 解析请求
			var statsReq post.GetPostsStatsRequest
			uriErr := c.ShouldBindUri(&statsReq)
			queryErr := c.ShouldBindQuery(&statsReq)

			// 检查基本绑定错误
			if tt.expectError && (tt.errorMsg == "min" || tt.errorMsg == "required") {
				if tt.errorMsg == "min" {
					// 对于min验证错误，可能同时有uri和query错误
					hasError := uriErr != nil || queryErr != nil
					assert.True(t, hasError)
					if uriErr != nil {
						assert.Contains(t, uriErr.Error(), tt.errorMsg)
					} else if queryErr != nil {
						assert.Contains(t, queryErr.Error(), tt.errorMsg)
					}
				} else if tt.errorMsg == "required" {
					// 对于required验证错误，可能同时有uri和query错误
					hasError := uriErr != nil || queryErr != nil
					assert.True(t, hasError)
					if uriErr != nil {
						assert.Contains(t, uriErr.Error(), tt.errorMsg)
					} else if queryErr != nil {
						assert.Contains(t, queryErr.Error(), tt.errorMsg)
					}
				}
				return
			}

			// 如果基本绑定成功，测试ParsePostIDs
			if uriErr == nil && queryErr == nil {
				postIDs, parseErr := statsReq.ParsePostIDs()

				if tt.expectError {
					assert.Error(t, parseErr)
					if tt.errorMsg != "" {
						assert.Contains(t, parseErr.Error(), tt.errorMsg)
					}
				} else {
					assert.NoError(t, parseErr)
					assert.Len(t, postIDs, tt.expectedLen)
					if tt.expectedLen > 0 {
						// 验证第一个ID是否正确解析
						expectedFirstID, _ := strconv.ParseInt("456", 10, 64)
						assert.Equal(t, expectedFirstID, postIDs[0])
					}
				}
			}
		})
	}
}

// TestGetPostsStatsRequestEdgeCases 测试获取文章统计信息请求的边界情况
func TestGetPostsStatsRequestEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("空字符串postIDs", func(t *testing.T) {
		req := post.GetPostsStatsRequest{PostIDs: ""}
		postIDs, err := req.ParsePostIDs()
		assert.NoError(t, err)
		assert.Empty(t, postIDs)
	})

	t.Run("只有逗号的postIDs", func(t *testing.T) {
		req := post.GetPostsStatsRequest{PostIDs: ",,,"}
		postIDs, err := req.ParsePostIDs()
		assert.NoError(t, err)
		assert.Empty(t, postIDs)
	})

	t.Run("包含空白项的postIDs", func(t *testing.T) {
		req := post.GetPostsStatsRequest{PostIDs: "123,,456,  ,789"}
		postIDs, err := req.ParsePostIDs()
		assert.NoError(t, err)
		assert.Len(t, postIDs, 3)
		assert.Equal(t, []int64{123, 456, 789}, postIDs)
	})

	t.Run("超大数字的postIDs", func(t *testing.T) {
		req := post.GetPostsStatsRequest{PostIDs: "9223372036854775807"} // int64最大值
		postIDs, err := req.ParsePostIDs()
		assert.NoError(t, err)
		assert.Len(t, postIDs, 1)
		assert.Equal(t, int64(9223372036854775807), postIDs[0])
	})

	t.Run("超出int64范围的数字", func(t *testing.T) {
		req := post.GetPostsStatsRequest{PostIDs: "9223372036854775808"} // 超出int64最大值
		_, err := req.ParsePostIDs()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid post id")
	})
}
