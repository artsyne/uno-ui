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
func Set(router *gin.Engine) {

	// 初始化 Views 模板
	setTemplate(router)

	// 初始化静态文件目录
	setStaticFs(router)

	// 初始化站点动态目录
	setPath(router)
}

// 通过rice box 加载 "views" 目录中的模版文件
func setTemplate(router *gin.Engine) {
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
	router.SetHTMLTemplate(tmpl)
}

// 设置静态文件目录，存放js、css，以及上传文件目录upload
func setStaticFs(router *gin.Engine) {
	//初始化静态文件
	//better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	router.StaticFS("/public", http.Dir(settings.PublicPath()))
	router.StaticFS("/upload", http.Dir(settings.PublicPath()))
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
}

// 设置站点router
func setPath(router *gin.Engine) {

	router.GET("/index", controllers.RenderIndex)

	// daybook
	daybook := router.Group("/daybook")
	{
		daybook.GET("/:spaceid", controllers.RenderDaybookIndex)
	}

	//表单渲染
	form := router.Group("/form")
	{
		form.GET("/daybook", controllers.RenderDaybookForm)
	}

	// diary
	diary := router.Group("/diary")
	{
		diary.GET("/", controllers.RenderDiary)
		diary.GET("/write", controllers.RenderDiaryWrite)
	}
	router.POST("/info", controllers.HandleInfo)

	// my
	router.GET("/my", controllers.RenderMyIndex)

	// passport
	router.GET("/passport/", controllers.RenderPassport)
	router.GET("/register", controllers.RenderRegister)
	router.POST("/register", controllers.HandleUserReg)
	router.GET("/login", controllers.RenderLogin)
	router.POST("/login", controllers.HandleLogin)
	router.POST("/loginout", controllers.HandleLoginOut)
	router.GET("/findpwd", controllers.RenderFindPwd)
	router.POST("/findpwd", controllers.HandleFindPwd)
	router.GET("/resetpwd", controllers.RenderResetPwd)
	router.POST("/resetpwd", controllers.HandleResetPwd)
	router.GET("/activateuser", controllers.RenderActivateUser)

	router.GET("/u/:uid", controllers.RenderMemberIndex)

	member := router.Group("/member")
	{
		member.GET("/", controllers.RenderMemberIndex)
	}

	// 日记本操作
	router.POST("/daybook", controllers.HandleSpace)
}
