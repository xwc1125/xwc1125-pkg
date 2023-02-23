// Package aesutil
//
// @author: xwc1125
package aesutil

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"
)

// 测试AES ECB 加密解密
func TestEncryptDecrypt(t *testing.T) {
	origData := []byte("Hello World") // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	log.Println("原文：", string(origData))

	log.Println("------------------ CBC模式 --------------------")
	encrypted := AesEncryptCBC(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesDecryptCBC(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	log.Println("------------------ ECB模式 --------------------")
	encrypted = AesEncryptECB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = AesDecryptECB(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	log.Println("------------------ CFB模式 --------------------")
	encrypted = AesEncryptCFB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = AesDecryptCFB(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	keys := "1234567890123456"
	enStr := "F0jppJSqb0ucBSpJSm4c2hH7d2Q3kDtKeHer3QRJEBGh+n2/RjgX2oBU+/gfXw6tN/2mbAjqzMYIf+ZEszi0PFqgCAbDAKyK+DENKrKW5uBb75SJeclntup12UhBVOgX"

	bytes, _ := base64.StdEncoding.DecodeString(enStr)
	decrypted = AesDecryptECB(bytes, []byte(keys))
	log.Println("解密结果1111：", string(decrypted))
}
