package render

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

// I18nManager 国际化管理器
type I18nManager struct {
	service  *RenderService
	messages map[string]map[string]string
	mu       sync.RWMutex
}

// NewI18nManager 创建国际化管理器
func NewI18nManager(service *RenderService) *I18nManager {
	return &I18nManager{
		service:  service,
		messages: make(map[string]map[string]string),
	}
}

// LoadMessages 加载国际化消息
func (i *I18nManager) LoadMessages(lang string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.messages[lang]; exists {
		return nil
	}

	// 构建消息文件路径
	messageFile := filepath.Join("website", "static", "messages", lang+".json")

	// 读取消息文件
	data, err := ioutil.ReadFile(messageFile)
	if err != nil {
		// 如果文件不存在，创建空的消息映射
		i.messages[lang] = make(map[string]string)
		return nil
	}

	// 解析JSON
	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return err
	}

	i.messages[lang] = messages
	return nil
}

// GetMessage 获取国际化消息
func (i *I18nManager) GetMessage(key string, args ...interface{}) string {
	// 默认语言
	lang := "en"

	// 如果第一个参数是字符串，则作为语言代码
	if len(args) > 0 {
		if langStr, ok := args[0].(string); ok {
			lang = langStr
		}
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	// 如果没有加载过该语言的消息，尝试加载
	if _, exists := i.messages[lang]; !exists {
		i.mu.RUnlock()
		i.LoadMessages(lang)
		i.mu.RLock()
	}

	// 查找消息
	if messages, exists := i.messages[lang]; exists {
		if message, found := messages[key]; found {
			return message
		}
	}

	// 尝试默认语言（英文）
	if lang != "en" {
		if messages, exists := i.messages["en"]; exists {
			if message, found := messages[key]; found {
				return message
			}
		}
	}

	// 返回键名
	return key
}

// GetMessageWithParams 获取带参数的国际化消息
func (i *I18nManager) GetMessageWithParams(key string, lang string, params map[string]string, defaultMessage ...string) string {
	message := i.GetMessage(key, lang)

	// 如果没有找到消息且提供了默认消息，使用默认消息
	if message == key && len(defaultMessage) > 0 {
		message = defaultMessage[0]
	}

	// 替换参数
	for paramKey, paramValue := range params {
		placeholder := "{{" + paramKey + "}}"
		message = strings.ReplaceAll(message, placeholder, paramValue)
	}

	return message
}

// GetMessageWithLang 获取指定语言的消息
func (i *I18nManager) GetMessageWithLang(key, lang string, args ...interface{}) string {
	return i.GetMessage(key, lang)
}

// HasMessage 检查是否存在消息
func (i *I18nManager) HasMessage(key string) bool {
	return i.GetMessage(key, "en") != key
}

// GetSupportedLanguages 获取支持的语言列表
func (i *I18nManager) GetSupportedLanguages() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	languages := make([]string, 0, len(i.messages))
	for lang := range i.messages {
		languages = append(languages, lang)
	}

	return languages
}

// RefreshMessages 刷新消息缓存
func (i *I18nManager) RefreshMessages() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.messages = make(map[string]map[string]string)
}

// GetUserLanguage 从用户信息中获取语言设置
func (i *I18nManager) GetUserLanguage(user RenderUser) string {
	if user == nil {
		return "en"
	}

	lang := user.GetLang()
	if lang == "" {
		return "en"
	}

	return lang
}

// SetMessage 设置消息（用于运行时动态设置）
func (i *I18nManager) SetMessage(key, message string) {
	i.SetMessageWithLang("en", key, message)
}

// SetMessageWithLang 设置指定语言的消息
func (i *I18nManager) SetMessageWithLang(lang, key, message string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.messages[lang]; !exists {
		i.messages[lang] = make(map[string]string)
	}

	i.messages[lang][key] = message
}
