package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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

func Test3Desc(t *testing.T) {
	key := "!@#$%^&*()" + string(utils.Rand(14, utils.RAND_KIND_ALL))
	t.Log(string(key))

	result, err := utils.TripleDesEncrypt("polaris@studygol", key)
	if err != nil {
		panic(err)
	}
	t.Log(result)
	if err != nil {
		panic(err)
	}
	origData, err := utils.TripleDesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	if string(origData) != "polaris@studygol" {
		t.Error("3Des加解密测试失败")
	}
	t.Log(string(origData))

}

func TestCreateRand(t *testing.T) {
	num := string(utils.Rand(14, utils.RAND_KIND_ALL))
	if len(num) != 14 {
		t.Error("创建随机数失败")
	}
	t.Log("创建随机数：" + num)
}

func TestSendEmail(t *testing.T) {

	user := "noreply@sinodiary.com"
	password := "HL90rmWKeJjS92"
	host := "smtpdm.aliyun.com:25"
	to := "biuren@gmail.com"

	subject := "使用Golang发送邮件"

	body := `
			<html>
			<body>
			<h3>
			"Test send to email"
			</h3>
			</body>
			</html>
			`
	fmt.Println("send email")
	err := utils.SendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		panic(err)
	} else {
		t.Log("Send mail success!")
	}
}

func TestGetJobsTasks(t *testing.T) {
	tasks, err := models.GetJobTasks(3, "email")
	if err != nil {
		panic(err)
	} else {
		str, _ := json.Marshal(tasks)
		t.Log(string(str))
	}
}

func TestAddValidCode(t *testing.T) {
	validCode := &models.ValidCode{}
	validCode.ValidCode = string(utils.Rand(32, utils.RAND_KIND_ALL))
	validCode.UserId = 35
	validCode.ValidTypeId = settings.VALID_CODE_ACTIVE_USER_EMAIL
	nowTime := utils.GetNowTime()
	validCode.CreateUTCTime = nowTime.UnixMilli()
	validCode.ExpireUTCTime = nowTime.AddDate(0, 0, 1).UnixMilli() //增加一天过期

	id, err := validCode.Insert()

	if err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Printf("valid code id= %d", id))
	}
}

func TestAddActivateUserEmailTask(t *testing.T) {
	user := &models.User{}

	user.Email = "kupage@gmail.com"
	user.UserId = 34
	user.UserName = "icefire"

	if err := services.AddActivateUserEmailTask(user); err != nil {
		t.Error(err)
	}
}

func TestCreateActiveUrl(t *testing.T) {
	//插入发送邮件队列
	config := settings.GetConfig()
	params := settings.GetParams()
	//激活链接
	url := fmt.Sprintf("%s/activiceuser?uid=%d&sign=%s", config.Domain, 10000, "wesdfsdfwerwerwer234234")
	t.Log(fmt.Sprintf(params["reg_email_body"], url, url))
}

func TestRedis(t *testing.T) {

	err := utils.SetSession("username", "kupage12ww")
	if err != nil {
		fmt.Println(err)
	}

	val, err := utils.GetSession("username")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(len(val))

	num, err := utils.DelSession("username")

	fmt.Println(num)
}

func TestGetSpacesByUserID(t *testing.T) {
	space, err := models.GetSpaceBySpaceID(6)
	fmt.Print(space)
	fmt.Print(err)
	fmt.Print(space.Title)
}

func TestGetSpaceByUserID(t *testing.T) {
	spaceArray, err := models.GetSpaceByUserID(34)
	fmt.Print(err)
	fmt.Print(len(spaceArray))
	fmt.Print(spaceArray)
}

func TestAddSpace(t *testing.T) {
	space := &models.Space{}
	space.Title = "测试日记"
	space.Remark = "这是备注"
	space.AccessPwd = "this is pwd"
	space.AccessURL = "http://accessurl"
	space.AddUTCTime = utils.GetNowTimeUnixMilli()
	space.AddUserID = 123
	space.AddUserName = "yasong"
	space.CoverMidURL = "http://covermidurl"
	space.CoverURL = "http://coverurl"
	space.IsShowComment = 1
	space.IsShowLike = 2
	space.IsShowView = 3
	space.LikeNum = 5
	space.MemberNum = 6
	space.Status = 1
	space.ViewNum = 7
	space.ViewSetting = 3

	id, err := space.Insert()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(id)
}

// 创建表对象code
func TestCreateCode(t *testing.T) {
	tables, err := models.GetTables()
	if err != nil {
		fmt.Println(err)
	}

	for _, tableObj := range tables {
		fmt.Println(PrintTableGolangStructCode(tableObj.TableName) + "\n\n")
	}
}

// 根据表明创建struct code 代码
func PrintTableGolangStructCode(tableName string) string {
	columns, err := models.GetTableColumns(tableName)

	if err != nil {
		fmt.Println(err)
	}
	//Struct Name, table name 格式为
	tableNameStr := ""
	for index, str := range strings.Split(tableName, "_") {
		if index != 0 {
			tableNameStr += strings.ToUpper(str[0:1]) + str[1:]
		}
	}
	outputCode := fmt.Sprintf("type %s struct {\n", tableNameStr)

	for _, col := range columns {
		fieldName := ""
		//循环转换列
		for _, str := range strings.Split(col.Field, "_") {

			tempStr := strings.ToUpper(str[0:1]) + str[1:]
			switch str {
			case "id":
				tempStr = "ID"
			case "ip":
				tempStr = "IP"
			case "url":
				tempStr = "URL"
			}

			fieldName += tempStr
		}
		//转换数据类型
		dataType := ""
		switch {
		case strings.Index(col.Type, "bigint") >= 0:
			dataType = "int64"
		case strings.Index(col.Type, "int(11)") >= 0:
			dataType = "int32"
		case strings.Index(col.Type, "tinyint(4)") >= 0:
			dataType = "int8"
		case col.Type == "datetime":
			dataType = "time.Time"
		default:
			dataType = "string"
		}

		outputCode += fmt.Sprintf("%s			%s		`form:\"%s\" json:\"%s\" db:\"%s\"` \n",
			fieldName, dataType, col.Field, col.Field, col.Field)
	}
	outputCode += fmt.Sprintf("}")

	return outputCode
}

func TestCreateFileName(t *testing.T) {
	filePath, fileName := utils.CreateUploadPathFileName()

	filePath = "upload/" + filePath
	fmt.Println(filePath + fileName)

	// 检查路径是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果路径不存在，则创建目录
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			t.Error(err)
		}
	}
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

//测试Config，反序列化输出值
func TestConfig(t *testing.T) {
	config := settings.GetConfig()

	// 将配置文件序列化为 JSON 格式的字符串
    configJSON, err := json.Marshal(config)

    if err != nil {
        panic(err)
    }
	 // 输出序列化后的字符串
	 fmt.Println(string(configJSON))
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