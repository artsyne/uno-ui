package tests

import (
	"encoding/json"
	"fmt"
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

// TestRand 测试 Rand 函数
func TestRand(t *testing.T) {
	size := 32
	kinds := []int{utils.RAND_KIND_NUM, utils.RAND_KIND_LOWER, utils.RAND_KIND_UPPER, utils.RAND_KIND_ALL}

	for _, kind := range kinds {
		result := utils.Rand(size, kind)
		t.Log(result)

		if len(result) != size {
			t.Errorf("Expected length of %d, but got %d", size, len(result))
		}

		for _, r := range result {
			if (kind == utils.RAND_KIND_NUM && (r < '0' || r > '9')) ||
				(kind == utils.RAND_KIND_LOWER && (r < 'a' || r > 'z')) ||
				(kind == utils.RAND_KIND_UPPER && (r < 'A' || r > 'Z')) ||
				(kind == utils.RAND_KIND_ALL && !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'))) {
				t.Errorf("Unexpected character: %c", r)
			}
		}
	}
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
	space, err := models.GetDiaryByDiaryID(6)
	fmt.Print(space)
	fmt.Print(err)
	fmt.Print(space.Title)
}

func TestAddSpace(t *testing.T) {
	space := &models.Diary{}
	space.Title = "测试日记"
	space.Remark = "这是备注"
	space.AccessPwd = "this is pwd"
	space.AccessURL = "http://accessurl"
	space.AddUTCTime = utils.GetNowTimeUnixMilli()
	space.AddUserID = 123
	space.AddUserName = "yasong"
	space.CoverMidURL = "http://covermidurl"
	space.CoverURL = "http://coverurl"
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

// 测试Config，反序列化输出值
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

// TestIsValidUserName 测试 IsValidUserName 函数
func TestIsValidUserName(t *testing.T) {
	// 定义测试用例
	cases := []struct {
		input string
		want  bool
	}{
		{"valid_user", true},                         // 有效的用户名
		{"invalid user", false},                      // 包含空格的用户名
		{"short", false},                             // 长度不足的用户名
		{"thisusernameiswaytoolongtobevalid", false}, // 长度过长的用户名
	}

	// 运行测试用例
	for _, c := range cases {
		got := utils.IsValidUserName(c.input)
		if got != c.want {
			t.Errorf("IsValidUserName(%q) == %v, want %v", c.input, got, c.want)
		}
	}
}

// TestIsValidPassword 测试 IsValidPassword 函数
func TestIsValidPassword(t *testing.T) {
	// 定义测试用例
	cases := []struct {
		input string
		want  bool
	}{
		{"validPassword1", true},                      // 有效的密码
		{"noDigit", false},                            // 不包含数字的密码
		{"short1", false},                             // 长度不足的密码
		{"thispasswordiswaytoolongtobevalid1", false}, // 长度过长的密码
		{"noLetter1", false},                          // 不包含字母的密码
	}

	// 运行测试用例
	for _, c := range cases {
		got := utils.IsValidPassword(c.input)
		if got != c.want {
			t.Errorf("IsValidPassword(%q) == %v, want %v", c.input, got, c.want)
		}
	}
}
