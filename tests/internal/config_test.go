package tests

import (
    "testing"

    "github.com/gin-gonic/gin"
    
    "github.com/biuren/hebe/internal/configs"
)

// 测试 GetConfig 方法
func TestGetConfig(t *testing.T) {
    // 设置 gin 模式
    gin.SetMode(gin.DebugMode)

    // 加载配置
    configs.LoadConfig()

    // 验证配置加载
    if configs.GetConfig() == nil {
        t.Fatal("Config should not be nil")
    }
    if configs.GetConfig().Domain != "http://127.0.0.1:8080" {
        t.Errorf("Expected domain to be 'localhost', got '%s'", configs.GetConfig().Domain)
    }
}

// 测试 GetParams 方法
func TestGetParams(t *testing.T) {
    // 设置 gin 模式
    gin.SetMode(gin.DebugMode)

    // 加载配置
    configs.LoadConfig()

    // 验证参数加载
    params := configs.GetParams()
    if len(params) == 0 {
        t.Fatal("Params should not be empty")
    }
    if params["reg_email_subject"] != "中国日记网注册激活邮件" {
        t.Errorf("Expected param1 to be 'value1', got '%s'", params["param1"])
    }
    // 
    t.Log(params["reg_email_body"])
}

// 测试 GetTimeZone 方法
func TestGetTimeZone(t *testing.T) {
    // 设置 gin 模式
    gin.SetMode(gin.DebugMode)

    // 加载配置
    configs.LoadConfig()

    // 验证时区加载
    timeZones := configs.GetTimeZone()
    if len(timeZones) == 0 {
        t.Fatal("Time zones should not be empty")
    }
    if timeZones[1] != "Africa/Algiers" {
        t.Errorf("Expected first time zone to be 'UTC', got '%s'", timeZones[0])
    }
}

// 测试 GetSupportLangs 方法
func TestGetSupportLangs(t *testing.T) {
    // 设置 gin 模式
    gin.SetMode(gin.DebugMode)

    // 加载配置
    configs.LoadConfig()

    // 验证支持语言加载
    supportLangs := configs.GetSupportLangs()
    if len(supportLangs) == 0 {
        t.Fatal("Support languages should not be empty")
    }
    if supportLangs[1].Code != "de" {
        t.Errorf("Expected first support language code to be 'en', got '%s'", supportLangs[0].Code)
    }
}