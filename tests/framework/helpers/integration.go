package helpers

import (
	htmlTemplate "html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/biuren/hebe/internal/app/templates"
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/gin-gonic/gin"
)

// LoadTestTemplates loads HTML templates for testing
func LoadTestTemplates(router *gin.Engine) {
	// Determine path to views
	// Try multiple levels up to find website/views
	basePath := "website/views"
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(basePath); err == nil {
			break
		}
		basePath = "../" + basePath
	}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// Fallback for direct execution from root
		basePath = "website/views"
	}

	tmpl := htmlTemplate.New("").Funcs(templates.RegisterAllFuncs())

	err := filepath.Walk(basePath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".html") {
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = tmpl.New(relPath).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(tmpl)
}

// MockAuthMiddleware to inject user into context
func MockAuthMiddleware(userMap map[string]*types.AuthUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetHeader("X-Test-User-ID")
		if uid != "" {
			if user, ok := userMap[uid]; ok {
				c.Set("IsLogin", true)
				c.Set("AuthUser", user)
			}
		}
		c.Next()
	}
}
