package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha256"
	"encoding/base64"
)

// Sha256 函数接受一个字符串作为参数，对其进行 SHA256 哈希，
// 然后对哈希结果进行 Base64 编码。函数返回编码后的字符串。
func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str)) //将字符串转换为[]bye
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// TripleDesEncrypt 函数接受一个待加密的字符串和一个密钥作为参数，
// 使用 3DES 算法对字符串进行加密，然后对加密结果进行 Base64 编码。
// 如果加密成功，函数返回编码后的字符串和 nil 错误；
func TripleDesEncrypt(enstr, key string) (string, error) {
	byteKey := []byte(key)
	block, err := des.NewTripleDESCipher(byteKey)
	if err != nil {
		return "", err
	}
	origData := []byte(enstr)
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, byteKey[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

// TripleDesDecrypt 函数接受一个已加密的字符串和一个密钥作为参数，
// 使用 3DES 算法对字符串进行解密。如果解密成功，函数返回解密后的字节切片和 nil 错误；
func TripleDesDecrypt(crypted, key string) ([]byte, error) {
	byteKey := []byte(key)
	block, err := des.NewTripleDESCipher(byteKey)
	if err != nil {
		return nil, err
	}
	byteCrypted, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, byteKey[:8])
	origData := make([]byte, len(byteCrypted))
	blockMode.CryptBlocks(origData, byteCrypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5Padding 函数接受一个字节切片和一个块大小作为参数，
// 对字节切片进行 PKCS5 填充。函数返回填充后的字节切片。
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding 函数接受一个已填充的字节切片作为参数，
// 对字节切片进行 PKCS5 去填充。函数返回去填充后的字节切片。
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
