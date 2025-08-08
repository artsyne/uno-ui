package view

import (
	"github.com/biuren/hebe/internal/utils/config"
)

// ===================== 配置服务 =====================

// GetSupportLangs 获取支持的语言列表
func (s *ViewService) GetSupportLangs() []config.SupportLang {
	return s.config.GetSupportLangs()
}

// GetTimeZones 获取时区列表
func (s *ViewService) GetTimeZones() []string {
	return s.config.GetTimeZone()
}

// getSiteName 获取网站名称
func (s *ViewService) getSiteName() string {
	if s.config != nil {
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			return name
		}
	}
	return "iDiary"
}