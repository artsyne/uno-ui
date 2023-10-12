package controllers

import (
	"fmt"
	"net/http"
	"path"
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

	//处理上传图片
	_, fileName := uploadImage(c, infoForm.CoverImgData)

	info := &models.Info{}
	var nowTime = helpers.GetNowTimeUnixMilli()
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
	info.AddUTCTime = nowTime
	info.AddUserIP = c.ClientIP()
	info.UpdateUTCTime = nowTime
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

// 上传文件图片
func uploadImage(c *gin.Context, imgBase64Data string) (filePath, fileName string) {

	//判断图片是否为空
	fileSize := len(imgBase64Data)
	if fileSize < 100 {
		return "", ""
	}

	// 获取上传文件配置信息
	config := system.GetConfig()

	//判断上传文件类型是否正确
	filetype := imgBase64Data[:strings.IndexByte(imgBase64Data, ';')]
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
	filePath = path.Join(config.Upload.ImgSavePath, filePath)
	fileName = fmt.Sprintf("%s/%s.jpeg", filePath, fileName)

	//保存图片
	err := helpers.SaveBase64File(imgBase64Data, filePath, fileName)
	if err != nil {
		globalErr.ResponseErr(c, globalErr.UPLOAD_FILE_ERR)
		return
	}

	return filePath, fileName
}
