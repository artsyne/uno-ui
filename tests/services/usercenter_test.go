package tests

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/services"
	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/utils"
)

// 初始化数据库连接
func init() {
	settings.LoadConfig()
	models.ConnectToDB()
}


// 测试时间函数
func TestUTCUnix(t *testing.T) {
	t.Log(utils.GetNowTime().UnixMicro())

	t.Log(utils.GetNowTimeUnixMilli())

	nowTime := utils.GetTimeFromUnixMilli(utils.GetNowTimeUnixMilli())

	t.Log(nowTime)

	t.Log(utils.GetTimeFromUnixMilli(1695711841053))
	t.Log(utils.GetTimeFromUnixMilli(1695798241053))

	t.Log(1695798241053 > utils.GetNowTimeUnixMilli())

}

// 根据timezone字符串获取时区，显示正确的本地时间
func TestGetTimeByTimeZone(t *testing.T) {
	utcTime := utils.GetNowTimeUnixMilli()
	t.Log(utcTime)

	location := "WET"
	localTime, _ := utils.GetLocalTimeFromUTC(utcTime, location)

	t.Log(localTime.Format("2006-01-02 15:04"))

	// 读取 timezone.json 文件
	file, err := os.Open("../config/timezone.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 解析 JSON 数据
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var timeZones []string
	if err := json.Unmarshal(data, &timeZones); err != nil {
		panic(err)
	}

	// 输出所有时区信息
	for _, tz := range timeZones {
		t.Log(tz)
	}
}

//测试获取用户语言
func TestGetUserLanguage(t *testing.T) {
	
	config := settings.GetConfig()
	
	//测试台湾语言
	if  services.GetSupportedLang("zh-tw", config.SupportLangs) != "zh-tw" {
		t.Error("获取用户语言zh-tw失败")
	}
	//测试中文
	if  services.GetSupportedLang("zh-cn", config.SupportLangs) != "zh-cn" {
		t.Error("获取用户语言zh-cn失败")
	}
	//测试英文
	if  services.GetSupportedLang("en-us", config.SupportLangs) != "en" {
		t.Error("获取用户语言en-us失败")
	}
	//测试日文	
	if  services.GetSupportedLang("Ja-jp", config.SupportLangs) != "ja" {
		t.Error("获取用户语言ja-jp失败")
	}
	//测试默认语言
	if  services.GetSupportedLang("pl-PL", config.SupportLangs) != "en" {
		t.Error("获取用户语言fr-fr失败")
	}
}