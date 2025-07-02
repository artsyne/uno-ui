package testhelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/module/user"
	"github.com/gin-gonic/gin"
)

// HTTPTestHelper HTTP测试辅助工具
type HTTPTestHelper struct {
	router *gin.Engine
}

// SetupHTTPTest 初始化HTTP测试环境
func SetupHTTPTest() *HTTPTestHelper {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	// 添加必要的中间件
	router.Use(user.SharedData())

	return &HTTPTestHelper{
		router: router,
	}
}

// GetRouter 获取路由器实例
func (h *HTTPTestHelper) GetRouter() *gin.Engine {
	return h.router
}

// GET 发送GET请求
func (h *HTTPTestHelper) GET(t *testing.T, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return h.makeRequest(t, "GET", path, nil, headers)
}

// POST 发送POST请求
func (h *HTTPTestHelper) POST(t *testing.T, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return h.makeRequest(t, "POST", path, body, headers)
}

// PUT 发送PUT请求
func (h *HTTPTestHelper) PUT(t *testing.T, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return h.makeRequest(t, "PUT", path, body, headers)
}

// DELETE 发送DELETE请求
func (h *HTTPTestHelper) DELETE(t *testing.T, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return h.makeRequest(t, "DELETE", path, nil, headers)
}

// PostForm 发送表单数据
func (h *HTTPTestHelper) PostForm(t *testing.T, path string, formData map[string]string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	form := url.Values{}
	for key, value := range formData {
		form.Set(key, value)
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	return h.makeRequest(t, "POST", path, strings.NewReader(form.Encode()), headers)
}

// makeRequest 构建并发送HTTP请求
func (h *HTTPTestHelper) makeRequest(t *testing.T, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody io.Reader

	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = strings.NewReader(v)
		case []byte:
			reqBody = bytes.NewReader(v)
		case io.Reader:
			reqBody = v
		default:
			// JSON序列化
			jsonData, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("序列化请求体失败: %v", err)
			}
			reqBody = bytes.NewReader(jsonData)

			// 设置Content-Type
			if headers == nil {
				headers = make(map[string]string)
			}
			if headers["Content-Type"] == "" {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)

	return w
}

// WithAuth 添加认证信息到请求头
func (h *HTTPTestHelper) WithAuth(t *testing.T, userID int64, username string) map[string]string {
	t.Helper()

	// 加密用户数据（需要获取配置中的密钥）
	// 这里需要根据实际的加密方式来实现
	// 暂时返回基本的认证头
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %d:%s", userID, username),
	}
}

// AssertStatus 断言HTTP状态码
func (h *HTTPTestHelper) AssertStatus(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()

	if w.Code != expectedStatus {
		t.Errorf("HTTP状态码不匹配: 期望 %d, 实际 %d", expectedStatus, w.Code)
	}
}

// AssertJSON 断言响应是JSON格式
func (h *HTTPTestHelper) AssertJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("响应不是JSON格式: Content-Type = %s", contentType)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fatalf("解析JSON响应失败: %v", err)
	}

	return data
}

// AssertJSONField 断言JSON字段值
func (h *HTTPTestHelper) AssertJSONField(t *testing.T, data map[string]interface{}, field string, expected interface{}) {
	t.Helper()

	actual, exists := data[field]
	if !exists {
		t.Errorf("JSON字段 '%s' 不存在", field)
		return
	}

	if actual != expected {
		t.Errorf("JSON字段 '%s' 值不匹配: 期望 %v, 实际 %v", field, expected, actual)
	}
}

// AssertContains 断言响应体包含指定内容
func (h *HTTPTestHelper) AssertContains(t *testing.T, w *httptest.ResponseRecorder, content string) {
	t.Helper()

	body := w.Body.String()
	if !strings.Contains(body, content) {
		t.Errorf("响应体不包含期望内容 '%s'", content)
	}
}

// AssertRedirect 断言重定向
func (h *HTTPTestHelper) AssertRedirect(t *testing.T, w *httptest.ResponseRecorder, expectedLocation string) {
	t.Helper()

	if w.Code != http.StatusFound && w.Code != http.StatusMovedPermanently && w.Code != http.StatusSeeOther {
		t.Errorf("响应不是重定向: 状态码 %d", w.Code)
		return
	}

	location := w.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("重定向位置不匹配: 期望 '%s', 实际 '%s'", expectedLocation, location)
	}
}

// AssertAPISuccess 断言API成功响应
func (h *HTTPTestHelper) AssertAPISuccess(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	h.AssertStatus(t, w, http.StatusOK)
	data := h.AssertJSON(t, w)
	h.AssertJSONField(t, data, "Status", float64(0)) // JSON unmarshaling converts numbers to float64

	return data
}

// AssertAPIError 断言API错误响应
func (h *HTTPTestHelper) AssertAPIError(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) map[string]interface{} {
	t.Helper()

	h.AssertStatus(t, w, expectedStatus)
	data := h.AssertJSON(t, w)

	// 检查是否有错误信息
	if _, exists := data["Message"]; !exists {
		t.Error("API错误响应缺少Message字段")
	}

	return data
}

// GetResponseData 获取响应中的Data字段
func (h *HTTPTestHelper) GetResponseData(t *testing.T, data map[string]interface{}) map[string]interface{} {
	t.Helper()

	dataField, exists := data["Data"]
	if !exists {
		t.Fatal("响应中缺少Data字段")
	}

	dataMap, ok := dataField.(map[string]interface{})
	if !ok {
		t.Fatal("Data字段不是对象类型")
	}

	return dataMap
}
