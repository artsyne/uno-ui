package controllers

import (
	"net/http"
	"strings"

	globalErr "github.com/biuren/hebe/config"
	"github.com/biuren/hebe/helpers"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/system"

	"github.com/gin-gonic/gin"
)

// renderDiaryWrite ： 渲染日记编辑
func RenderDiaryWrite(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "写日记"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "diary/write_form", h)
}

// 发布日记
func PostInfo(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	// 序列化请求对象
	infoForm := &models.InfoForm{}
	err := c.Bind(infoForm)
	if err != nil {
		globalErr.ResponseErr(c, globalErr.PARSE_INPUT_PARAM_ERR)
		return
	}

	//上传图片
	filePath := ""
	fileName := ""
	if fileSize := len(infoForm.CoverImgData); fileSize > 100 {
		// 获取上传文件配置信息
		config := system.GetConfig()

		//判断上传文件类型是否正确
		filetype := infoForm.CoverImgData[:strings.IndexByte(infoForm.CoverImgData, ';')]
		if !strings.Contains(config.Upload.ImgType, filetype) {
			globalErr.ResponseErr(c, globalErr.UPLOAD_FILE_TYPE_ERR)
			return
		}

		//判断文件大小
		if fileSize > config.Upload.ImgMaxSize {
			globalErr.ResponseErr(c, globalErr.UPLOAD_FILE_TOO_LARGE)
			return
		}

		//路径格式为：config配置路径 + /年/月/日/uid.jpg
		filePath, fileName = helpers.CreateUploadPathFileName()
		filePath = config.Upload.ImgSavePath + filePath
		fileName = filePath + fileName + ".jpeg"

		//保存图片
		err := helpers.SaveBase64Img(infoForm.CoverImgData, filePath, fileName)
		if err != nil {
			globalErr.ResponseErr(c, globalErr.UPLOAD_FILE_ERR)
			return
		}
	}

	info := &models.Info{}
	var nowTime = helpers.GetNowTime()
	//赋值info对象
	info.SpaceID = infoForm.SpaceID
	info.AccessURL = fileName
	info.CoverURL = fileName
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
		globalErr.ResponseErrWithDesc(c, globalErr.DB_OPERATION_ERR,
			globalErr.EnumErrDesc[globalErr.DB_OPERATION_ERR]+","+err.Error())
		return
	}

	//创建返回对象
	h := helpers.ResponseH(c)
	h["Data"] = info
	h["Form"] = infoForm
	h["infoId"] = infoID
	c.JSON(http.StatusOK, h)
}
