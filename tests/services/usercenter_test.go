package tests

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/modules/user"
	"github.com/biuren/hebe/pkg"
)

// 初始化数据库连接
func init() {
	configs.InitConfig()
	models.ConnectToDB()
}

// 测试时间函数
func TestUTCUnix(t *testing.T) {
	t.Log(pkg.GetNowTime().UnixMicro())

	t.Log(pkg.GetNowTimeUnixMilli())

	nowTime := pkg.GetTimeFromUnixMilli(pkg.GetNowTimeUnixMilli())

	t.Log(nowTime)

	t.Log(pkg.GetTimeFromUnixMilli(1695711841053))
	t.Log(pkg.GetTimeFromUnixMilli(1695798241053))

	t.Log(1695798241053 > pkg.GetNowTimeUnixMilli())

}

// 根据timezone字符串获取时区，显示正确的本地时间
func TestGetTimeByTimeZone(t *testing.T) {
	utcTime := pkg.GetNowTimeUnixMilli()
	t.Log(utcTime)

	location := "WET"
	localTime, _ := pkg.GetLocalTimeFromUTC(utcTime, location)

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

// 测试获取用户语言
func TestGetUserLanguage(t *testing.T) {
	userService := user.NewService()

	//测试台湾语言
	if userService.GetSupportedLang("zh-tw") != "zh-tw" {
		t.Error("获取用户语言zh-tw失败")
	}
	//测试中文
	if userService.GetSupportedLang("zh-cn") != "zh-cn" {
		t.Error("获取用户语言zh-cn失败")
	}
	//测试英文
	if userService.GetSupportedLang("en-us") != "en" {
		t.Error("获取用户语言en-us失败")
	}
}

func Backward(s []string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			yield(strings.ToUpper(s[i]))
		}
	}
}

// 测试golang 1.23 新特性
func TestGolang123(t *testing.T) {
	//测试golang 1.23 新特性
	t.Log("测试golang 1.23 新特性")

	s := []string{"hello", "world", "golang", "GO to hall"}
	backwardFunc := Backward(s)
	backwardFunc(func(str string) bool {
		t.Log(str)
		return true
	})
}
