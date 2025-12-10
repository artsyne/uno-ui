package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"github.com/biuren/hebe/internal/app/templates"
)

// Mock functions to avoid initializing the full i18n stack if not strictly necessary for syntax check
// or we can use the real one. Let's try to minimal mock first.
// If the error is about a missing function, we will see it.
// If the error is syntax in the HTML, we will see it.

func main() {
	basePath := "/Users/yasong.cy/go/src/idiary.ai/hebe/website/views"
	
	// Register all funcs
	funcs := templates.RegisterAllFuncs()
	tmpl := template.New("").Funcs(funcs)

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

			fmt.Printf("Parsing %s...\n", relPath)
			// Parse
			_, err = tmpl.New(relPath).Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", relPath, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("TEMPLATE ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All templates parsed successfully!")
}
