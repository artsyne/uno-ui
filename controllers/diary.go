package controllers

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	gobalerr "chinadiary.com/website/config"
	"chinadiary.com/website/helpers"
	"chinadiary.com/website/models"
	"chinadiary.com/website/system"

	"github.com/gin-gonic/gin"
)

// renderDiaryWrite ： 渲染日记编辑
func RenderDiaryWrite(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "写日记"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "diary/write_form", h)
}

//发布日记
func PostInfo(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	// 序列化请求对象
	infoForm := &models.InfoForm{}
	err := c.Bind(infoForm)
	if err != nil {
		gobalerr.ResponseErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
		return
	}

	//上传图片
	if len(infoForm.CoverImgData) > 100 {
		// 获取上传文件配置信息
		config := system.GetConfig()

		//判断上传文件类型是否正确
		filetype := infoForm.CoverImgData[:strings.IndexByte(infoForm.CoverImgData, ';')]
		if !strings.Contains(config.Upload.ImgType, filetype) {
			gobalerr.ResponseErr(c, gobalerr.UPLOAD_FILE_TYPE_ERR)
			return
		}

		//图片存储路径包含文件名，不带扩展名
		//路径格式为：路径 + /年/月/日/uid-imgfile
		filename := config.Upload.ImgSavePath + "xxx.jpg"

		//去掉,号前内容，"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYA"
		b64data := infoForm.CoverImgData[strings.IndexByte(infoForm.CoverImgData, ',')+1:]

		//图片解码
		imgByte, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			gobalerr.ResponseErr(c, gobalerr.UPLOAD_FILE_ERR)
			return
		}
		err = ioutil.WriteFile(filename, imgByte, 0666)
		if err != nil {
			gobalerr.ResponseErr(c, gobalerr.UPLOAD_FILE_ERR)
			return
		}
	}

	info := &models.Info{}
	var nowTime = helpers.GetNowTime()
	//赋值info对象
	info.SpaceID = infoForm.SpaceID
	info.AccessURL = ""
	info.CoverURL = ""
	info.CoverMidURL = ""
	info.CoverStyle = infoForm.CoverStyle
	info.IsTop = 0
	info.Title = infoForm.Title
	info.Summary = ""
	info.Body = infoForm.Body
	info.ViewNum = 0
	info.LikeNum = 0
	info.CommentNum = 0
	info.AddUserID = userData.UserId
	info.AddUserName = userData.UserName
	info.AddTime = nowTime
	info.AddUserIP = c.ClientIP()
	info.UpdateTime = nowTime
	info.Status = 1

	//插入数据库
	infoID, err := info.Insert()
	if err != nil {
		gobalerr.ResponseErrWithDesc(c, gobalerr.DB_OPERATION_ERR,
			gobalerr.EnumErrDesc[gobalerr.DB_OPERATION_ERR]+","+err.Error())
		return
	}

	//创建返回对象
	h := helpers.ResponseH(c)
	h["Data"] = info
	h["Form"] = infoForm
	h["infoId"] = infoID
	c.JSON(http.StatusOK, h)
}
