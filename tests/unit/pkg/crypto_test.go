package pkg_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"testing"
	"time"

	"github.com/biuren/hebe/pkg"
	core "github.com/biuren/hebe/tests/framework/core"
)



// TestSha256 测试 Sha256 函数
func TestSha256(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}
	
	str := "Hello, world!"
	hash := sha256.New()
	hash.Write([]byte(str))
	expected := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if pkg.Sha256(str) != expected {
		t.Errorf("Sha256 did not return the expected result")
	}
}


func TestAESEncryptDecrypt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
    // 测试用例
    testCases := []struct {
        plainText string
        key       string
    }{
        {"Hello, world!", "1234567890123456"}, // 16 字节，对应 128 位
        {"Go is awesome!", "abcdefghijklmnopqrstuvwx"}, // 24 字节，对应 192 位
        {"GitHub Copilot", "2GzVloukJl9C7y44O740o59e3sD8xTzr"}, // 32 字节，对应 256 位
    }

    for _, tc := range testCases {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
        // 使用 AESEncrypt 函数加密明文
        cipherText, err :=pkg.AESEncrypt(tc.plainText, tc.key)
		t.Log(cipherText)
        if err != nil {
            t.Errorf("AESEncrypt(%q, %q) returned error: %v", tc.plainText, tc.key, err)
        }

        // 使用 AESDecrypt 函数解密密文
        decryptedText, err := pkg.AESDecrypt(cipherText, tc.key)
		t.Log(string(decryptedText))
        if err != nil {
            t.Errorf("AESDecrypt(%q, %q) returned error: %v", cipherText, tc.key, err)
        }

        // 检查解密后的文本是否与原始明文相同
        if string(decryptedText) != tc.plainText {
            t.Errorf("Decrypted text %q does not match original plain text %q", decryptedText, tc.plainText)
        }
    }
}

// TestPasswordHashingPerformance 比较 SHA256+Salt 和 Bcrypt 的性能
// TestBcryptHash 测试 Bcrypt 哈希功能
func TestBcryptHash(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}
	
	password := "testPassword123!"

	// 测试 Bcrypt 哈希
	hash, err := pkg.BcryptHash(password)
	if err != nil {
		t.Fatalf("Bcrypt 哈希失败: %v", err)
	}

	// 验证哈希值不为空
	if hash == "" {
		t.Error("Bcrypt 哈希值为空")
	}

	// 验证哈希值长度（bcrypt 哈希值通常为 60 字符）
	if len(hash) != 60 {
		t.Errorf("Bcrypt 哈希值长度不正确，期望 60，实际 %d", len(hash))
	}

	// 验证密码
	if !pkg.BcryptVerify(password, hash) {
		t.Error("Bcrypt 密码验证失败")
	}

	// 验证错误密码
	if pkg.BcryptVerify("wrongPassword", hash) {
		t.Error("Bcrypt 应该拒绝错误密码")
	}

	t.Logf("Bcrypt 哈希值: %s", hash)
	t.Logf("Bcrypt 哈希值长度: %d 字符", len(hash))
}

// TestBcryptPerformance 测试 Bcrypt 性能
func TestBcryptPerformance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}
	
	password := "testPassword123!"

	// 测试 Bcrypt 哈希性能（默认 cost）
	start := time.Now()
	hash, err := pkg.BcryptHash(password)
	if err != nil {
		t.Fatalf("Bcrypt 哈希失败: %v", err)
	}
	hashDuration := time.Since(start)
	t.Logf("Bcrypt 哈希时间 (默认 Cost): %v", hashDuration)

	// 测试 Bcrypt 验证性能
	start = time.Now()
	verifyResult := pkg.BcryptVerify(password, hash)
	verifyDuration := time.Since(start)
	t.Logf("Bcrypt 验证时间: %v", verifyDuration)

	// 验证结果正确性
	if !verifyResult {
		t.Error("Bcrypt 验证失败")
	}

	t.Logf("Bcrypt 总时间: %v", hashDuration+verifyDuration)
}