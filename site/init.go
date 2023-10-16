package site

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/biuren/hebe/helpers"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/system"
	"github.com/biuren/hebe/task"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 初始化站点，链接数据，设置cookie
// 返回默认gin 的router对象
func InitSite() *gin.Engine {
	//链接数据库
	ConnectToDB()
	logrus.Info("数据库连接成功！")

	//设置cookie默认值
	config := system.GetConfig()
	helpers.SetCookieDefVal(config.Cookie.Domain, "/", config.Cookie.ExpireDays, config.Cookie.HttpOnly)
	logrus.Info("Cookie默认值设置成功！")

	//启动定时任务
	task.Run()
	logrus.Info("启动定时任务成功！")

	//初始化默认gin router
	router := gin.Default()

	return router
}

// 初始化Router，加载静态文件，上下文共享数据，中间件
func SetRouter(router *gin.Engine) {
	//初始化模板
	setTemplate(router)
	logrus.Info("初始化模板成功！")

	//初始化静态文件
	//better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	router.StaticFS("/public", http.Dir(system.PublicPath()))
	router.StaticFS("/upload", http.Dir(system.PublicPath()))
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	logrus.Info("初始化静态文件成功！")

	//加载上下文共享数据，如登录数据
	router.Use(SharedData())
}

// http 中间件
// 获取用户登录信息
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取cookie
		cookieEnUserData := helpers.ReadCookie(c.Request, "ud")
		//获取是否有cookie
		if len(cookieEnUserData) < 1 {
			return
		}

		//解密cookie
		byteUserData, err := helpers.TripleDesDecrypt(cookieEnUserData, system.GetConfig().Cookie.SecretKey)
		if err != nil || byteUserData == nil {
			return
		}

		//反序列化为对象
		var userData models.Login
		err = json.Unmarshal(byteUserData, &userData)
		if err != nil {
			return
		}

		//判断认证信息

		//写入上下文
		c.Set("UserData", userData)
		c.Set("IsLogin", true)
		c.Next()
	}
}

// 初始化数据库链接
func ConnectToDB() {
	config := system.GetConfig()
	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local", config.Database.User, config.Database.Password,
		config.Database.Host, config.Database.Name)
	models.SetDB(connection)
}

// 通过rice box 加载 "views" 目录中的模版文件
func setTemplate(router *gin.Engine) {
	box := rice.MustFindBox("../views")
	//注册模版函数
	tmpl := template.New("").Funcs(template.FuncMap{})

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
