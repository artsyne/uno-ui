package localization

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/biuren/hebe/internal/core/config"
)

// Loader 文件加载器
type Loader struct {
	cfg        *config.Config
	normalizer *Normalizer
}

// NewLoader 创建文件加载器
func NewLoader(cfg *config.Config, normalizer *Normalizer) *Loader {
	return &Loader{
		cfg:        cfg,
		normalizer: normalizer,
	}
}

// LoadAllMessages 加载所有语言的消息到全局缓存
func (l *Loader) LoadAllMessages() error {
	// 使用配置中的 LocalesDir
	localesDir := l.cfg.App.LocalesDir

	// 确保目录存在
	if _, err := os.Stat(localesDir); os.IsNotExist(err) {
		return fmt.Errorf("locales目录不存在: %s", localesDir)
	}

	loadedCount := 0
	supportedCodes := l.normalizer.GetSupportedLanguageCodes()

	for _, langCode := range supportedCodes {
		if err := l.loadLanguageMessages(localesDir, langCode); err != nil {
			slog.Error("加载语言消息失败", "language", langCode, "error", err)
			// 继续加载其他语言，不中断
			continue
		}
		loadedCount++
	}

	if loadedCount == 0 {
		return fmt.Errorf("没有成功加载任何消息")
	}

	slog.Info("消息加载完成", "成功加载语言数", loadedCount, "总支持语言数", len(supportedCodes))
	return nil
}

// loadLanguageMessages 加载指定语言的消息
func (l *Loader) loadLanguageMessages(localesDir, langCode string) error {
	// 处理语言代码映射（zh-cn -> zh等）
	normalizedLangCode := l.normalizer.NormalizeLocaleForFile(langCode)
	langFile := filepath.Join(localesDir, fmt.Sprintf("%s.json", normalizedLangCode))

	// 如果规范化后的文件不存在，尝试原始代码
	if _, err := os.Stat(langFile); os.IsNotExist(err) {
		langFile = filepath.Join(localesDir, fmt.Sprintf("%s.json", langCode))
	}

	// 检查文件是否存在
	if _, err := os.Stat(langFile); os.IsNotExist(err) {
		slog.Warn("语言文件不存在", "language", langCode, "file", langFile)
		return fmt.Errorf("语言文件不存在: %s", langFile)
	}

	// 读取文件内容
	data, err := os.ReadFile(langFile)
	if err != nil {
		return fmt.Errorf("读取语言文件失败 %s: %w", langFile, err)
	}

	// 解析 JSON
	var langMessages map[string]interface{}
	if err := json.Unmarshal(data, &langMessages); err != nil {
		return fmt.Errorf("解析语言文件失败 %s: %w", langFile, err)
	}

	// 扁平化并存储消息
	l.flattenAndStore(langMessages, langCode, "")

	slog.Debug("成功加载语言消息", "language", langCode, "file", langFile, "消息数", len(langMessages))
	return nil
}

// flattenAndStore 扁平化嵌套的消息结构并存储到全局缓存
func (l *Loader) flattenAndStore(data map[string]interface{}, langCode, prefix string) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			// 直接存储字符串消息
			cacheKey := BuildCacheKey(langCode, fullKey)
			SetMessage(cacheKey, v)
		case map[string]interface{}:
			// 递归处理嵌套对象
			l.flattenAndStore(v, langCode, fullKey)
		default:
			// 其他类型转为字符串
			cacheKey := BuildCacheKey(langCode, fullKey)
			SetMessage(cacheKey, fmt.Sprintf("%v", v))
		}
	}
}
