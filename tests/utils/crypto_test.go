package tests

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/biuren/hebe/utils"
)

// TestSha256 测试 Sha256 函数
func TestSha256(t *testing.T) {
	str := "Hello, world!"
	hash := sha256.New()
	hash.Write([]byte(str))
	expected := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if utils.Sha256(str) != expected {
		t.Errorf("Sha256 did not return the expected result")
	}
}

// TestTripleDesEncrypt 测试 TripleDesEncrypt 函数
func TestTripleDesEncrypt(t *testing.T) {
	enstr := "Hello, world!Hello, world!Hello, world!"
	key := "EGCF7q9RF8SIYM73G732r1XK"
	enStr, err := utils.TripleDesEncrypt(enstr, key)
	t.Log(enStr)
	if err != nil {
		t.Errorf("TripleDesEncrypt returned an error: %v", err)
	}
}

// TestTripleDesDecrypt 测试 TripleDesDecrypt 函数
func TestTripleDesDecrypt(t *testing.T) {
	key := "EGCF7q9RF8SIYM73G732r1XK"
	crypted := "PIpR2NdKFNvfIv/SW74mmft8M0KJIT+mu9yIRRQUXScgBnky7Dnmgg==" // 这应该是一个已经使用 3DES 加密并且 Base64 编码后的字符串
	deStr, err := utils.TripleDesDecrypt(crypted, key)
	t.Log(string(deStr))
	if err != nil {
		t.Errorf("TripleDesDecrypt returned an error: %v", err)
	}
}

// 测试3Des加解密
func Test3Desc(t *testing.T) {
	key := string(utils.Rand(24, utils.RAND_KIND_ALL))
	t.Log(len(key))
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