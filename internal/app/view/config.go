package view

import (
	"github.com/biuren/hebe/internal/utils/config"
)

// ===================== 配置服务 =====================

// GetSupportLangs 获取支持的语言列表
// 通过配置服务获取系统支持的所有语言配置
//
// 返回:
//   - []config.SupportLang: 支持的语言配置列表
func (s *ViewService) GetSupportLangs() []config.SupportLang {
	return s.config.GetSupportLangs()
}

// GetTimeZones 获取时区列表
// 通过配置服务获取系统支持的所有时区列表
//
// 返回:
//   - []string: 时区字符串列表，如"Asia/Shanghai"、"UTC"等
func (s *ViewService) GetTimeZones() []string {
	return s.config.GetTimeZone()
}

// getSiteName 获取网站名称
// 从配置参数中获取网站名称，如果未配置则返回默认名称
//
// 返回:
//   - string: 网站名称，默认为"iDiary"
func (s *ViewService) getSiteName() string {
	if s.config != nil {
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			return name
		}
	}
	return "iDiary"
}