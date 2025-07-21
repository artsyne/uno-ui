package templates

import "html/template"

// RegisterAllFuncs 注册所有模板函数
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
