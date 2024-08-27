package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha256"
	"encoding/base64"
	"errors"
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

// AESEncrypt 加密明文数据，返回 base64 编码的密文。
// 参数 plainText 是要加密的明文数据。
// 参数 key 是用于加密的密钥，它应该是一个 16、24、32位字符串
// 返回的结果是 base64 编码的密文。
func AESEncrypt(plainText string, key string) (string, error) {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }

    // 对明文数据进行 PKCS7 填充
    origData := PKCS5Padding([]byte(plainText), block.BlockSize())
    
    // 创建一个新的分组模式 cipher.BlockMode
    blockMode := cipher.NewCBCEncrypter(block, []byte(key)[:aes.BlockSize])

    cipherText := make([]byte, len(origData))
    blockMode.CryptBlocks(cipherText, origData)

    return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AESDecrypt 解密 base64 编码的密文，返回明文数据。
// 参数 cipherTextStr 是要解密的 base64 编码的密文。
// 参数 key 是用于解密的密钥，它应该是一个 16、24、32位字符串。
// 返回的结果是解密后的明文数据。
func AESDecrypt(cipherTextStr string, key string) ([]byte, error) {
    cipherText, err := base64.StdEncoding.DecodeString(cipherTextStr)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return nil, err
    }

    if len(cipherText) < aes.BlockSize {
        return nil, errors.New("cipherText too short")
    }

    // 创建一个新的分组模式 cipher.BlockMode
    blockMode := cipher.NewCBCDecrypter(block, []byte(key)[:aes.BlockSize])

    plainText := make([]byte, len(cipherText))
    blockMode.CryptBlocks(plainText, cipherText)

    // 移除 PKCS7 填充
    plainText = PKCS5UnPadding(plainText)

    return plainText, nil
}