package templates_test

import (
	"testing"

	"github.com/biuren/hebe/internal/app/templates"
	"github.com/stretchr/testify/assert"
)

func TestGetInitial(t *testing.T) {
	funcs := templates.GetUIFuncs()
	getInitialFunc, ok := funcs["getInitial"]
	assert.True(t, ok, "getInitial should be registered")

	getInitial, ok := getInitialFunc.(func(string) string)
	assert.True(t, ok, "getInitial should be func(string) string")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"English Name", "Alice", "A"},
		{"Lowercase English", "bob", "B"},
		{"Chinese Name", "张三", "张"},
		{"Empty", "", ""},
		{"Special Char", "@handle", "@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, getInitial(tt.input))
		})
	}
}
