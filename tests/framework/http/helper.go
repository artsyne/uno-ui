package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	framework "github.com/biuren/hebe/tests/framework/types"
)

// HTTPTestHelperImpl HTTP测试助手实现
type HTTPTestHelperImpl struct {
	router *gin.Engine
}

// NewHTTPTestHelper 创建HTTP测试助手
func NewHTTPTestHelper() *HTTPTestHelperImpl {
	return &HTTPTestHelperImpl{}
}

// SetupRouter 设置路由
func (h *HTTPTestHelperImpl) SetupRouter() error {
	h.router = gin.New()
	return nil
}

// request 通用请求方法
func (h *HTTPTestHelperImpl) request(method, path string, body interface{}, headers map[string]string) *framework.HTTPResponse {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return &framework.HTTPResponse{StatusCode: http.StatusInternalServerError}
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)

	return &framework.HTTPResponse{
		StatusCode: w.Code,
		Body:       w.Body.Bytes(),
		Headers:    map[string]string{"Content-Type": w.Header().Get("Content-Type")},
	}
}

// GET 发送GET请求
func (h *HTTPTestHelperImpl) GET(path string, headers map[string]string) *framework.HTTPResponse {
	return h.request("GET", path, nil, headers)
}

// POST 发送POST请求
func (h *HTTPTestHelperImpl) POST(path string, body interface{}, headers map[string]string) *framework.HTTPResponse {
	return h.request("POST", path, body, headers)
}

// PUT 发送PUT请求
func (h *HTTPTestHelperImpl) PUT(path string, body interface{}, headers map[string]string) *framework.HTTPResponse {
	return h.request("PUT", path, body, headers)
}

// DELETE 发送DELETE请求
func (h *HTTPTestHelperImpl) DELETE(path string, headers map[string]string) *framework.HTTPResponse {
	return h.request("DELETE", path, nil, headers)
}