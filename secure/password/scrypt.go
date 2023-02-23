// Package password
//
// @author: xwc1125
package password

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	pwHashBytes = 64
)

// GenerateSalt 生成盐
func GenerateSalt() (salt string, err error) {
	buf := make([]byte, pwHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", buf), nil
}

// GeneratePasswordByScrypt 生成passHash
func GeneratePasswordByScrypt(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, pwHashBytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h), nil
}
