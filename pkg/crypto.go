package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// Sha256 函数接受一个字符串作为参数，对其进行 SHA256 哈希，
// 然后对哈希结果进行 Base64 编码。函数返回编码后的字符串。
func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str)) //将字符串转换为[]bye
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// AESEncrypt 加密明文数据，返回 base64 编码的密文。
// 参数 plainText 是要加密的明文数据。
// 参数 key 是用于加密的密钥，它应该是一个 16、24、32位字符串
// 返回的结果是 base64 编码的密文。
func AESEncrypt(plainText string, key string) (string, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AESDecrypt 解密 base64 编码的密文，返回明文数据。
// 参数 cipherTextStr 是要解密的 base64 编码的密文。
// 参数 key 是用于解密的密钥，它应该是一个 16、24、32位字符串。
// 返回的结果是解密后的明文数据。
func AESDecrypt(cipherTextStr string, key string) ([]byte, error) {
	keyBytes := []byte(key)
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("cipherText too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// BcryptHash 使用 bcrypt 算法对密码进行哈希处理。
// 参数 password 是要哈希的明文密码。
// 返回 bcrypt 哈希值和可能的错误。
func BcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// BcryptVerify 验证明文密码与 bcrypt 哈希值是否匹配。
// 参数 password 是明文密码。
// 参数 hash 是 bcrypt 哈希值（其中已包含成本信息）。
// 返回验证是否成功的布尔值。
func BcryptVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}