package routes

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/biuren/hebe/controllers"
	"github.com/biuren/hebe/settings"
	"github.com/gin-gonic/gin"
)

// 初始化Router，加载静态文件，上下文共享数据，中间件
func Set(server *gin.Engine) {

	// 初始化 Views 模板
	setTemplate(server)

	// 初始化静态文件目录
	setStaticFs(server)

	// 初始化站点动态目录
	setPath(server)
}

// 通过rice box 加载 "views" 目录中的模版文件
func setTemplate(server *gin.Engine) {
	box := rice.MustFindBox("../views")
	//注册模版函数
	tmpl := template.New("").Funcs(template.FuncMap{
		//TODO：加载模版函数，后期需要完善
		"formatTime": func(utc int64, layout string) string {
			return box.Time().UTC().Format(layout)
		},
	})

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			var err error
			tmpl, err = tmpl.Parse(box.MustString(path))
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := box.Walk("", fn)
	if err != nil {
		panic(err)
	}
	server.SetHTMLTemplate(tmpl)
}

// 设置静态文件目录，存放js、css，以及上传文件目录upload
func setStaticFs(server *gin.Engine) {
	//初始化静态文件
	//better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	server.StaticFS("/public", http.Dir(settings.PublicPath()))
	server.StaticFS("/static", http.Dir(settings.PublicPath()))
	server.StaticFile("/favicon.ico", "./public/favicon.ico")
}

// 设置站点router
func setPath(server *gin.Engine) {

	server.GET("/index", controllers.RenderIndex)

	// daybook
	daybook := server.Group("/daybook")
	{
		daybook.GET("/:spaceid", controllers.RenderDaybookIndex)
	}

	//表单渲染
	form := server.Group("/form")
	{
		form.GET("/daybook", controllers.RenderDaybookForm)
	}

	// diary
	diary := server.Group("/diary")
	{
		diary.GET("/", controllers.RenderDiary)
		diary.GET("/write", controllers.RenderDiaryWrite)
	}
	server.POST("/info", controllers.HandleInfo)

	// my
	server.GET("/my", controllers.RenderMyIndex)

	// passport
	server.GET("/passport/", controllers.RenderPassport)
	server.GET("/register", controllers.RenderRegister)
	server.POST("/register", controllers.HandleUserReg)
	server.GET("/login", controllers.RenderLogin)
	server.POST("/login", controllers.HandleLogin)
	server.POST("/loginout", controllers.HandleLoginOut)
	server.GET("/findpwd", controllers.RenderFindPwd)
	server.POST("/findpwd", controllers.HandleFindPwd)
	server.GET("/resetpwd", controllers.RenderResetPwd)
	server.POST("/resetpwd", controllers.HandleResetPwd)
	server.GET("/activateuser", controllers.RenderActivateUser)

	server.GET("/u/:uid", controllers.RenderMemberIndex)

	member := server.Group("/member")
	{
		member.GET("/", controllers.RenderMemberIndex)
	}

	// 日记本操作
	server.POST("/daybook", controllers.HandleSpace)
}
