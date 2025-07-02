package configs_test

import (
	"testing"

	"github.com/biuren/hebe/internal/core/config"
)

func TestConfigInit(t *testing.T) {
	// 测试配置初始化
	err := config.InitConfig()
	if err != nil {
		t.Errorf("配置初始化失败: %v", err)
	}

	// 获取配置
	config := config.GetConfig()
	if config == nil {
		t.Error("获取配置失败，config为nil")
	}
}

func TestConstants(t *testing.T) {
	// 测试一些基本常量
	if config.DEFAULT_USER_GROUP <= 0 {
		t.Error("DEFAULT_USER_GROUP应该大于0")
	}

	if config.USER_STATUS_WAIT_ACTIVE < 0 {
		t.Error("USER_STATUS_WAIT_ACTIVE应该大于等于0")
	}
}
