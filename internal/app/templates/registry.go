package templates

import "html/template"

// RegisterAllFuncs 注册所有模板函数
// 统一注册所有分类的模板函数，包括本地化、认证、UI和时间相关的函数
//
// 返回:
//   - template.FuncMap: 包含所有模板函数的映射
//
// 包含的函数分类:
//   - 本地化函数: getMessage、getMessageWithDefault等
//   - 认证函数: isLogin、authUser、getUserName等
//   - UI函数: 用户界面相关的辅助函数
//   - 时间函数: 时间格式化和处理相关函数
func RegisterAllFuncs() template.FuncMap {
	funcs := make(template.FuncMap)

	// 合并所有分类的模板函数
	for name, fn := range GetLocalizationFuncs() {
		funcs[name] = fn
	}

	for name, fn := range GetAuthFuncs() {
		funcs[name] = fn
	}

	for name, fn := range GetUIFuncs() {
		funcs[name] = fn
	}

	for name, fn := range GetTimeFuncs() {
		funcs[name] = fn
	}

	return funcs
}
