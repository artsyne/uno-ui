package localization

import (
	"strings"

	"github.com/biuren/hebe/internal/core/config"
)

// Normalizer 语言代码规范化器
type Normalizer struct {
	supportedLangs map[string]config.SupportLang
	defaultLang    string
}

// NewNormalizer 创建语言代码规范化器
func NewNormalizer(cfg *config.Config) *Normalizer {
	supportedLangs := make(map[string]config.SupportLang)
	for _, lang := range cfg.SupportLangs {
		supportedLangs[lang.Code] = lang
	}

	return &Normalizer{
		supportedLangs: supportedLangs,
		defaultLang:    cfg.App.DefaultLang,
	}
}

// NormalizeLocale 规范化语言代码（基于配置中的支持语言）
func (n *Normalizer) NormalizeLocale(locale string) string {
	if locale == "" {
		return n.defaultLang
	}

	// 统一转为小写
	locale = strings.ToLower(locale)

	// 直接检查是否为支持的语言
	if _, exists := n.supportedLangs[locale]; exists {
		return locale
	}

	// 处理语言代码变体（如 zh-CN -> zh-cn）
	for supportedCode := range n.supportedLangs {
		if strings.EqualFold(locale, supportedCode) {
			return supportedCode
		}

		// 处理主语言匹配（如 zh-Hans -> zh-cn）
		if strings.HasPrefix(locale, strings.Split(supportedCode, "-")[0]) {
			return supportedCode
		}
	}

	// 处理常见的语言代码映射
	normalizedMap := map[string]string{
		"zh":      "zh-cn",
		"zh-hans": "zh-cn",
		"zh-hant": "zh-tw",
		"zh-sg":   "zh-cn",
		"zh-hk":   "zh-tw",
		"zh-mo":   "zh-tw",
	}

	if mapped, exists := normalizedMap[locale]; exists {
		if _, supported := n.supportedLangs[mapped]; supported {
			return mapped
		}
	}

	return n.defaultLang
}

// NormalizeLocaleForFile 规范化语言代码用于文件名
func (n *Normalizer) NormalizeLocaleForFile(locale string) string {
	// 将 zh-cn 转换为 zh 等，用于查找文件
	switch locale {
	case "zh-cn", "zh-tw":
		return "zh"
	default:
		return strings.Split(locale, "-")[0] // 取主语言部分
	}
}

// IsLanguageSupported 检查语言是否支持
func (n *Normalizer) IsLanguageSupported(locale string) bool {
	normalized := n.NormalizeLocale(locale)
	_, exists := n.supportedLangs[normalized]
	return exists
}

// GetLanguageInfo 获取语言信息
func (n *Normalizer) GetLanguageInfo(locale string) (config.SupportLang, bool) {
	normalized := n.NormalizeLocale(locale)
	lang, exists := n.supportedLangs[normalized]
	return lang, exists
}

// GetSupportedLanguages 获取支持的语言列表
func (n *Normalizer) GetSupportedLanguages() []config.SupportLang {
	langs := make([]config.SupportLang, 0, len(n.supportedLangs))
	for _, lang := range n.supportedLangs {
		langs = append(langs, lang)
	}
	return langs
}

// GetSupportedLanguageCodes 获取支持的语言代码列表
func (n *Normalizer) GetSupportedLanguageCodes() []string {
	codes := make([]string, 0, len(n.supportedLangs))
	for code := range n.supportedLangs {
		codes = append(codes, code)
	}
	return codes
}

// GetDefaultLanguage 获取默认语言
func (n *Normalizer) GetDefaultLanguage() string {
	return n.defaultLang
}
