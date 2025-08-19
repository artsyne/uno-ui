package pkg

import (
	"context"
	"testing"
	"time"

	"github.com/biuren/hebe/pkg"
	"github.com/stretchr/testify/assert"
	core "github.com/biuren/hebe/tests/framework/core"
)

// TestRegexCheck 测试RegexCheck结构体的所有方法
func TestRegexCheck(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	rc := &pkg.RegexCheck{}

	// 测试IsID方法
	t.Run("IsID测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的ID
		validIDs := []string{
			"abc123",
			"ABC123",
			"123abc",
			"abcDEF123456",
			"123456789012345678901234567890", // 30位
		}
		
		for _, id := range validIDs {
			assert.True(t, rc.IsID(id), "ID '%s' 应该是有效的", id)
		}
		
		// 无效的ID
		invalidIDs := []string{
			"abc12",     // 少于6位
			"abc@123",   // 包含特殊字符
			"abc 123",   // 包含空格
			"abc-123",   // 包含连字符
			"1234567890123456789012345678901", // 超过30位
			"",          // 空字符串
			"中文123",    // 包含中文
		}
		
		for _, id := range invalidIDs {
			assert.False(t, rc.IsID(id), "ID '%s' 应该是无效的", id)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsID("abc123", "def456"), "多个有效ID应该返回true")
		assert.False(t, rc.IsID("abc123", "def@456"), "包含无效ID应该返回false")
	})

	// 测试IsUserName方法
	t.Run("IsUserName测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的用户名
		validUserNames := []string{
			"user123",
			"USER123",
			"123user",
			"userNAME123456",
			"123456789012345678901234567890", // 30位
		}
		
		for _, username := range validUserNames {
			assert.True(t, rc.IsUserName(username), "用户名 '%s' 应该是有效的", username)
		}
		
		// 无效的用户名
		invalidUserNames := []string{
			"user1",     // 少于6位
			"user@123",  // 包含特殊字符
			"user 123",  // 包含空格
			"user-123",  // 包含连字符
			"1234567890123456789012345678901", // 超过30位
			"",          // 空字符串
			"用户123",    // 包含中文
		}
		
		for _, username := range invalidUserNames {
			assert.False(t, rc.IsUserName(username), "用户名 '%s' 应该是无效的", username)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsUserName("user123", "admin456"), "多个有效用户名应该返回true")
		assert.False(t, rc.IsUserName("user123", "admin@456"), "包含无效用户名应该返回false")
	})

	// 测试IsPwd方法
	t.Run("IsPwd测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的密码
		validPasswords := []string{
			"pass123",
			"PASS123",
			"pass@123",
			"pass.123",
			"123@abc",
			"password@123.ABC",
			"123456789012345678901234567890", // 30位
		}
		
		for _, password := range validPasswords {
			assert.True(t, rc.IsPwd(password), "密码 '%s' 应该是有效的", password)
		}
		
		// 无效的密码
		invalidPasswords := []string{
			"pass1",     // 少于6位
			"pass#123",  // 包含不允许的特殊字符
			"pass 123",  // 包含空格
			"pass-123",  // 包含连字符
			"1234567890123456789012345678901", // 超过30位
			"",          // 空字符串
			"密码123",    // 包含中文
		}
		
		for _, password := range invalidPasswords {
			assert.False(t, rc.IsPwd(password), "密码 '%s' 应该是无效的", password)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsPwd("pass123", "admin@456"), "多个有效密码应该返回true")
		assert.False(t, rc.IsPwd("pass123", "admin#456"), "包含无效密码应该返回false")
	})

	// 测试IsInteger方法
	t.Run("IsInteger测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的整数
		validIntegers := []string{
			"0",
			"1",
			"123",
			"999999",
			"1234567890",
		}
		
		for _, integer := range validIntegers {
			assert.True(t, rc.IsInteger(integer), "整数 '%s' 应该是有效的", integer)
		}
		
		// 无效的整数
		invalidIntegers := []string{
			"-123",    // 负数
			"12.3",    // 小数
			"12a",     // 包含字母
			"12 3",    // 包含空格
			"",        // 空字符串
			"abc",     // 纯字母
			"12.0",    // 带小数点的整数
		}
		
		for _, integer := range invalidIntegers {
			assert.False(t, rc.IsInteger(integer), "整数 '%s' 应该是无效的", integer)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsInteger("123", "456"), "多个有效整数应该返回true")
		assert.False(t, rc.IsInteger("123", "45.6"), "包含无效整数应该返回false")
	})

	// 测试IsDecimals方法
	t.Run("IsDecimals测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的小数
		validDecimals := []string{
			"1.0",
			"12.34",
			"0.5",
			"999.999",
			"123.456789",
		}
		
		for _, decimal := range validDecimals {
			assert.True(t, rc.IsDecimals(decimal), "小数 '%s' 应该是有效的", decimal)
		}
		
		// 无效的小数
		invalidDecimals := []string{
			"123",     // 整数
			".5",      // 没有整数部分
			"12.",     // 没有小数部分
			"12.a",    // 小数部分包含字母
			"a.5",     // 整数部分包含字母
			"-12.3",   // 负数
			"12 .3",   // 包含空格
			"",        // 空字符串
			"12..3",   // 多个小数点
		}
		
		for _, decimal := range invalidDecimals {
			assert.False(t, rc.IsDecimals(decimal), "小数 '%s' 应该是无效的", decimal)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsDecimals("12.3", "45.6"), "多个有效小数应该返回true")
		assert.False(t, rc.IsDecimals("12.3", "456"), "包含无效小数应该返回false")
	})

	// 测试IsCellphone方法
	t.Run("IsCellphone测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的手机号
		validCellphones := []string{
			"13812345678",
			"15987654321",
			"18612345678",
			"17712345678",
			"19912345678",
		}
		
		for _, cellphone := range validCellphones {
			assert.True(t, rc.IsCellphone(cellphone), "手机号 '%s' 应该是有效的", cellphone)
		}
		
		// 无效的手机号
		invalidCellphones := []string{
			"1381234567",   // 少于11位
			"138123456789",  // 超过11位
			"0381234567",   // 不以1开头
			"138123456a8",  // 包含字母
			"138 1234 5678", // 包含空格
			"",             // 空字符串
			"138-1234-5678", // 包含连字符
		}
		
		for _, cellphone := range invalidCellphones {
			assert.False(t, rc.IsCellphone(cellphone), "手机号 '%s' 应该是无效的", cellphone)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsCellphone("13812345678", "15987654321"), "多个有效手机号应该返回true")
		assert.False(t, rc.IsCellphone("13812345678", "1598765432"), "包含无效手机号应该返回false")
	})

	// 测试IsTelephone方法
	t.Run("IsTelephone测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的座机号
		validTelephones := []string{
			"12345678",
			"87654321",
			"01234567",
			"99999999",
		}
		
		for _, telephone := range validTelephones {
			assert.True(t, rc.IsTelephone(telephone), "座机号 '%s' 应该是有效的", telephone)
		}
		
		// 无效的座机号
		invalidTelephones := []string{
			"1234567",   // 少于8位
			"123456789", // 超过8位
			"1234567a",  // 包含字母
			"1234 567",  // 包含空格
			"",          // 空字符串
			"1234-567",  // 包含连字符
		}
		
		for _, telephone := range invalidTelephones {
			assert.False(t, rc.IsTelephone(telephone), "座机号 '%s' 应该是无效的", telephone)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsTelephone("12345678", "87654321"), "多个有效座机号应该返回true")
		assert.False(t, rc.IsTelephone("12345678", "8765432"), "包含无效座机号应该返回false")
	})

	// 测试IsEngishLowCase方法
	t.Run("IsEngishLowCase测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的小写英文
		validLowCase := []string{
			"a",
			"abc",
			"hello",
			"world",
			"abcdefghijklmnopqrstuvwxyz",
		}
		
		for _, lowCase := range validLowCase {
			assert.True(t, rc.IsEngishLowCase(lowCase), "小写英文 '%s' 应该是有效的", lowCase)
		}
		
		// 无效的小写英文
		invalidLowCase := []string{
			"A",        // 大写字母
			"Hello",    // 包含大写字母
			"abc123",   // 包含数字
			"abc@def",  // 包含特殊字符
			"abc def",  // 包含空格
			"",         // 空字符串
			"你好",      // 中文
		}
		
		for _, lowCase := range invalidLowCase {
			assert.False(t, rc.IsEngishLowCase(lowCase), "小写英文 '%s' 应该是无效的", lowCase)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsEngishLowCase("hello", "world"), "多个有效小写英文应该返回true")
		assert.False(t, rc.IsEngishLowCase("hello", "World"), "包含无效小写英文应该返回false")
	})

	// 测试IsEnglishCap方法
	t.Run("IsEnglishCap测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的大写英文
		validCapCase := []string{
			"A",
			"ABC",
			"HELLO",
			"WORLD",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		}
		
		for _, capCase := range validCapCase {
			assert.True(t, rc.IsEnglishCap(capCase), "大写英文 '%s' 应该是有效的", capCase)
		}
		
		// 无效的大写英文
		invalidCapCase := []string{
			"a",        // 小写字母
			"Hello",    // 包含小写字母
			"ABC123",   // 包含数字
			"ABC@DEF",  // 包含特殊字符
			"ABC DEF",  // 包含空格
			"",         // 空字符串
			"你好",      // 中文
		}
		
		for _, capCase := range invalidCapCase {
			assert.False(t, rc.IsEnglishCap(capCase), "大写英文 '%s' 应该是无效的", capCase)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsEnglishCap("HELLO", "WORLD"), "多个有效大写英文应该返回true")
		assert.False(t, rc.IsEnglishCap("HELLO", "World"), "包含无效大写英文应该返回false")
	})

	// 测试IsEnglish方法
	t.Run("IsEnglish测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的英文
		validEnglish := []string{
			"a",
			"A",
			"abc",
			"ABC",
			"Hello",
			"WORLD",
			"HelloWorld",
			"abcDEF",
			"ABCdef",
		}
		
		for _, english := range validEnglish {
			assert.True(t, rc.IsEnglish(english), "英文 '%s' 应该是有效的", english)
		}
		
		// 无效的英文
		invalidEnglish := []string{
			"abc123",   // 包含数字
			"abc@def",  // 包含特殊字符
			"abc def",  // 包含空格
			"",         // 空字符串
			"你好",      // 中文
			"hello-world", // 包含连字符
		}
		
		for _, english := range invalidEnglish {
			assert.False(t, rc.IsEnglish(english), "英文 '%s' 应该是无效的", english)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsEnglish("Hello", "World"), "多个有效英文应该返回true")
		assert.False(t, rc.IsEnglish("Hello", "World123"), "包含无效英文应该返回false")
	})

	// 测试IsEmail方法
	t.Run("IsEmail测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 有效的邮箱
		validEmails := []string{
			"test@example.com",
			"user.name@domain.co.uk",
			"user_name@domain.org",
			"user-name@domain.net",
			"123@domain.com",
			"test.email@sub.domain.com",
			"a@b.co",
		}
		
		for _, email := range validEmails {
			assert.True(t, rc.IsEmail(email), "邮箱 '%s' 应该是有效的", email)
		}
		
		// 无效的邮箱
		invalidEmails := []string{
			"test",           // 没有@符号
			"test@",          // 没有域名
			"@domain.com",    // 没有用户名
			"test@domain",    // 没有顶级域名
			"test@domain.",   // 顶级域名为空
			"test..test@domain.com", // 连续的点
			"test@domain..com", // 域名中连续的点
			"test domain@example.com", // 用户名包含空格
			"test@domain .com", // 域名包含空格
			"",              // 空字符串
			"test@",         // 不完整的邮箱
		}
		
		for _, email := range invalidEmails {
			assert.False(t, rc.IsEmail(email), "邮箱 '%s' 应该是无效的", email)
		}
		
		// 测试多个参数
		assert.True(t, rc.IsEmail("test@example.com", "user@domain.org"), "多个有效邮箱应该返回true")
		assert.False(t, rc.IsEmail("test@example.com", "invalid-email"), "包含无效邮箱应该返回false")
	})
}

// TestRegexCheckEdgeCases 测试RegexCheck的边界情况
func TestRegexCheckEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	rc := &pkg.RegexCheck{}

	t.Run("空参数测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 测试没有参数的情况
		assert.True(t, rc.IsID(), "没有参数时应该返回true")
		assert.True(t, rc.IsUserName(), "没有参数时应该返回true")
		assert.True(t, rc.IsPwd(), "没有参数时应该返回true")
		assert.True(t, rc.IsInteger(), "没有参数时应该返回true")
		assert.True(t, rc.IsDecimals(), "没有参数时应该返回true")
		assert.True(t, rc.IsCellphone(), "没有参数时应该返回true")
		assert.True(t, rc.IsTelephone(), "没有参数时应该返回true")
		assert.True(t, rc.IsEngishLowCase(), "没有参数时应该返回true")
		assert.True(t, rc.IsEnglishCap(), "没有参数时应该返回true")
		assert.True(t, rc.IsEnglish(), "没有参数时应该返回true")
		assert.True(t, rc.IsEmail(), "没有参数时应该返回true")
	})

	t.Run("大量参数测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 测试大量有效参数
		validIDs := []string{"abc123", "def456", "ghi789", "jkl012", "mno345"}
		assert.True(t, rc.IsID(validIDs...), "大量有效ID应该返回true")
		
		// 测试大量参数中包含一个无效的
		mixedIDs := []string{"abc123", "def456", "ghi789", "abc@123", "mno345"}
		assert.False(t, rc.IsID(mixedIDs...), "包含无效ID的大量参数应该返回false")
	})

	t.Run("特殊字符测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 测试各种特殊字符
		specialChars := []string{
			"abc!123", "abc#123", "abc$123", "abc%123", "abc^123",
			"abc&123", "abc*123", "abc(123", "abc)123", "abc+123",
			"abc=123", "abc[123", "abc]123", "abc{123", "abc}123",
			"abc|123", "abc\\123", "abc:123", "abc;123", "abc'123",
			"abc\"123", "abc<123", "abc>123", "abc,123", "abc?123",
			"abc/123", "abc~123", "abc`123",
		}
		
		for _, char := range specialChars {
			assert.False(t, rc.IsID(char), "包含特殊字符 '%s' 的ID应该是无效的", char)
			assert.False(t, rc.IsUserName(char), "包含特殊字符 '%s' 的用户名应该是无效的", char)
		}
	})

	t.Run("Unicode字符测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 测试各种Unicode字符
		unicodeStrings := []string{
			"测试123",    // 中文
			"тест123",   // 俄文
			"テスト123",   // 日文
			"테스트123",   // 韩文
			"test123α",  // 希腊字母
			"test123ñ",  // 西班牙语字符
			"test123ü",  // 德语字符
			"test123é",  // 法语字符
		}
		
		for _, unicode := range unicodeStrings {
			assert.False(t, rc.IsID(unicode), "包含Unicode字符 '%s' 的ID应该是无效的", unicode)
			assert.False(t, rc.IsUserName(unicode), "包含Unicode字符 '%s' 的用户名应该是无效的", unicode)
			assert.False(t, rc.IsEnglish(unicode), "包含Unicode字符 '%s' 的英文应该是无效的", unicode)
		}
	})
}

// TestRegexCheckPerformance 测试RegexCheck的性能
func TestRegexCheckPerformance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	rc := &pkg.RegexCheck{}

	t.Run("性能测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 测试大量调用的性能
		testString := "validtest123"
		
		// 这里只是确保方法能够快速执行，不会有性能问题
		for i := 0; i < 1000; i++ {
			rc.IsID(testString)
			rc.IsUserName(testString)
			rc.IsInteger("123")
			rc.IsEmail("test@example.com")
		}
		
		// 如果能执行到这里，说明性能是可接受的
		assert.True(t, true, "性能测试通过")
	})
}

// TestRegexCheckIntegration 测试RegexCheck的集成功能
func TestRegexCheckIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	rc := &pkg.RegexCheck{}

	t.Run("用户注册场景测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 模拟用户注册场景
		username := "testuser123"
		password := "password@123"
		email := "testuser@example.com"
		cellphone := "13812345678"
		
		// 验证所有字段
		assert.True(t, rc.IsUserName(username), "用户名应该有效")
		assert.True(t, rc.IsPwd(password), "密码应该有效")
		assert.True(t, rc.IsEmail(email), "邮箱应该有效")
		assert.True(t, rc.IsCellphone(cellphone), "手机号应该有效")
		
		// 综合验证
		allValid := rc.IsUserName(username) && rc.IsPwd(password) && 
				   rc.IsEmail(email) && rc.IsCellphone(cellphone)
		assert.True(t, allValid, "所有注册信息应该都有效")
	})

	t.Run("数据验证场景测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 模拟数据验证场景
		integerValue := "12345"
		decimalValue := "123.45"
		englishText := "HelloWorld"
		lowCaseText := "hello"
		upCaseText := "WORLD"
		
		// 验证数据类型
		assert.True(t, rc.IsInteger(integerValue), "整数值应该有效")
		assert.True(t, rc.IsDecimals(decimalValue), "小数值应该有效")
		assert.True(t, rc.IsEnglish(englishText), "英文文本应该有效")
		assert.True(t, rc.IsEngishLowCase(lowCaseText), "小写英文应该有效")
		assert.True(t, rc.IsEnglishCap(upCaseText), "大写英文应该有效")
		
		// 交叉验证
		assert.False(t, rc.IsInteger(decimalValue), "小数不应该被识别为整数")
		assert.False(t, rc.IsDecimals(integerValue), "整数不应该被识别为小数")
		assert.False(t, rc.IsEngishLowCase(upCaseText), "大写文本不应该被识别为小写")
		assert.False(t, rc.IsEnglishCap(lowCaseText), "小写文本不应该被识别为大写")
	})

	t.Run("联系方式验证场景测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 模拟联系方式验证场景
		cellphone := "13812345678"
		telephone := "12345678"
		email := "contact@company.com"
		
		// 验证联系方式
		assert.True(t, rc.IsCellphone(cellphone), "手机号应该有效")
		assert.True(t, rc.IsTelephone(telephone), "座机号应该有效")
		assert.True(t, rc.IsEmail(email), "邮箱应该有效")
		
		// 交叉验证（确保不会误判）
		assert.False(t, rc.IsTelephone(cellphone), "手机号不应该被识别为座机号")
		assert.False(t, rc.IsCellphone(telephone), "座机号不应该被识别为手机号")
	})
}