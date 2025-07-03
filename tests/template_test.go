package tests

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/app/render"
	"github.com/biuren/hebe/internal/core/config"
)

func TestTemplateLoading(t *testing.T) {
	// 获取项目根目录（从tests目录回到根目录）
	cwd, _ := os.Getwd()
	projectRoot := filepath.Dir(cwd) // 从tests目录回到项目根目录
	basePath := filepath.Join(projectRoot, "website/views")

	// 创建渲染服务（已合并template、i18n功能）
	cfg := config.GetConfig()
	renderService := render.NewRenderService(cfg)

	// 注册所有模版函数
	tmpl := template.New("").Funcs(renderService.RegisterFuncs())

	// 加载所有 .html 文件
	err := loadTemplates(tmpl, basePath)
	if err != nil {
		t.Fatal("Template loading error:", err)
	}

	log.Println("Templates loaded successfully!")
}

// loadTemplates 递归遍历目录，加载所有 .html 文件
func loadTemplates(tmpl *template.Template, basePath string) error {
	return filepath.Walk(basePath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".html") {
			// 使用相对路径作为模板名称，避免命名冲突
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}

			// 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			log.Printf("Loading template: %s", relPath)

			// 使用相对路径作为模板名称解析
			_, err = tmpl.New(relPath).Parse(string(content))
			if err != nil {
				log.Printf("Error parsing template %s: %v", relPath, err)
				return err
			}
		}
		return nil
	})
}
