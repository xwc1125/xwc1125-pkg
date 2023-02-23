// Package otp
//
// @author: xwc1125
package otp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

// EncryptTOTP 加密totp
func EncryptTOTP(password []byte, salt []byte, totp []byte) ([]byte, error) {
	dk := pbkdf2.Key(password, salt, 10000, 32, sha3.New256)
	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(totp))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], totp)

	return []byte(base64.URLEncoding.EncodeToString(ciphertext)), nil
}

// DecryptTOTP 解密totp
func DecryptTOTP(password []byte, salt []byte, totp string) (string, error) {
	dk := pbkdf2.Key(password, salt, 10000, 32, sha3.New256)

	ciphertext, _ := base64.URLEncoding.DecodeString(totp)
	block, err := aes.NewCipher(dk)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
