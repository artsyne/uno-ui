package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha256"
	"encoding/base64"
)

//字符串做sha256后base64编码
func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str)) //将字符串转换为[]bye
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// 3DES加密，加密结果做base64编码
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

// 3DES解密
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

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
