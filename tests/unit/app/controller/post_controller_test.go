package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	postDomain "github.com/biuren/hebe/internal/domain/post"
	post "github.com/biuren/hebe/internal/app/controller/post"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/user"
)



// TestNewPostController 测试控制器创建
func TestNewPostController(t *testing.T) {
	cfg := &config.Config{}
	mockPostSvc := &postDomain.Service{}
	mockDiarySvc := &diary.Service{}
	mockAuth := &user.AuthHelper{}
	mockViewSvc := &view.ViewService{}
	mockUserSvc := &user.Service{}

	// 创建控制器
	controller := post.NewPostController(
		cfg,
		mockPostSvc,
		mockDiarySvc,
		mockAuth,
		mockViewSvc,
		mockUserSvc,
	)

	// 验证控制器不为空
	assert.NotNil(t, controller)
}

// TestDeletePost 测试删除日记请求解析
func TestDeletePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody string
		contentType string
		expectError bool
		errorMsg    string
		expectedID  int64
	}{
		{
			name:        "有效的删除请求",
			requestBody: `{"post_id": 123, "diary_id": 456}`,
			contentType: "application/json",
			expectError: false,
			errorMsg:    "",
			expectedID:  123,
		},
		{
			name:        "post_id为0的请求",
			requestBody: `{"post_id": 0, "diary_id": 456}`,
			contentType: "application/json",
			expectError: true,
			errorMsg:    "required",
			expectedID:  0,
		},
		{
			name:        "负数post_id的请求",
			requestBody: `{"post_id": -1, "diary_id": 456}`,
			contentType: "application/json",
			expectError: true,
			errorMsg:    "min",
			expectedID:  -1,
		},
		{
			name:        "缺少post_id的请求",
			requestBody: `{"diary_id": 456}`,
			contentType: "application/json",
			expectError: true,
			errorMsg:    "required",
			expectedID:  0,
		},
		{
			name:        "无效JSON格式",
			requestBody: `{"post_id": }`,
			contentType: "application/json",
			expectError: true,
			errorMsg:    "invalid character",
			expectedID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建HTTP请求
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/post/delete", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", tt.contentType)

			// 创建Gin上下文
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 测试请求解析逻辑
			var deleteReq post.DeletePostRequest
			err := c.ShouldBindJSON(&deleteReq)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, deleteReq.PostID)
			}
		})
	}
}

// TestDeletePostFormData 测试表单数据删除日记
func TestDeletePostFormData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/delete", strings.NewReader("post_id=123&diary_id=456"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	var deleteReq post.DeletePostRequest
	err := c.ShouldBind(&deleteReq)

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
			name:        "post_id为0",
			postID:      "0",
			expectError: true,
			errorMsg:    "required",
		},
		{
			name:        "负数post_id",
			postID:      "-1",
			expectError: true,
			errorMsg:    "min",
		},
		{
			name:        "非数字post_id",
			postID:      "abc",
			expectError: true,
			errorMsg:    "cannot unmarshal string",
		},
		{
			name:        "空post_id",
			postID:      "",
			expectError: true,
			errorMsg:    "cannot unmarshal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			requestBody := `{"post_id": ` + tt.postID + `, "diary_id": 456}`
			if tt.postID == "abc" || tt.postID == "" {
				requestBody = `{"post_id": "` + tt.postID + `", "diary_id": 456}`
			}

			req, _ := http.NewRequest("POST", "/post/delete", bytes.NewBufferString(requestBody))
			req.Header.Set("Content-Type", "application/json")

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			var deleteReq post.DeletePostRequest
			err := c.ShouldBindJSON(&deleteReq)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				postID, _ := strconv.ParseInt(tt.postID, 10, 64)
				assert.Equal(t, postID, deleteReq.PostID)
				assert.Equal(t, int64(456), deleteReq.DiaryID)
			}
		})
	}
}