package middlewares

import (
    "html/template"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/render"
)

// CustomHTMLRender 自定义 HTML 渲染器
type CustomHTMLRender struct {
    templates *template.Template
    GlobalVars gin.H
}

// NewCustomHTMLRender 创建一个新的自定义渲染器实例
func NewCustomHTMLRender(templates *template.Template, globalVars gin.H) *CustomHTMLRender {
    return &CustomHTMLRender{
        templates: templates,
        GlobalVars: globalVars,
    }
}

// Instance 实现 render.HTMLRender 接口
func (r *CustomHTMLRender) Instance(name string, data interface{}) render.Render {
    if data == nil {
        data = gin.H{}
    }

    // 将全局变量合并到 data 中
    if ginData, ok := data.(gin.H); ok {
        for k, v := range r.GlobalVars {
            ginData[k] = v
        }
    }

    return render.HTML{
        Template: r.templates,
        Name:     name,
        Data:     data,
    }
}

// Render 覆盖默认的 Render 方法，自动添加全局变量
func (r *CustomHTMLRender) Render(w http.ResponseWriter, name string, data interface{}) error {
    return r.Instance(name, data).Render(w)
}