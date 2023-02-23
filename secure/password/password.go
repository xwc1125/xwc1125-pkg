// Package password
//
// @author: xwc1125
package password

import (
	"github.com/chain5j/chain5j-pkg/crypto/scrypt"
	"github.com/chain5j/chain5j-pkg/crypto/signature/prime256v1"
	"github.com/pborman/uuid"
	"github.com/xwc1125/xwc1125-pkg/secure/shift"
	"golang.org/x/crypto/bcrypt"
)

var (
	BcryptCost = bcrypt.DefaultCost
	Base       = "www.xwc1125.com"
	BaseKeyLen = 16
)

type Key struct {
	Id         string
	PrivateKey []byte
}

// GeneratePassword 生成密码
// plaintextPwd 明文密码
// 对密码进行加密，获得密码hash
func GeneratePassword(plaintextPwd string) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(plaintextPwd), BcryptCost)
}

// CheckPassword 验证密码的正确性
func CheckPassword(hashPwd, plaintextPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(plaintextPwd))
	return err == nil
}

// EncryptKey 加密密码
func EncryptKey(key Key, auth string) ([]byte, error) {
	return scrypt.EncryptKey(&scrypt.Key{
		Id:         uuid.UUID(key.Id),
		PrivateKey: key.PrivateKey,
	}, getKey(auth), scrypt.StandardScryptN, scrypt.StandardScryptP)
}

// DecryptKey 解密
func DecryptKey(data []byte, auth string) (*Key, error) {
	key, err := scrypt.DecryptKey(data, getKey(auth))
	if err != nil {
		return nil, err
	}
	return &Key{
		Id:         key.Id.String(),
		PrivateKey: key.PrivateKey,
	}, nil
}

func getKey(auth string) string {
	salt := auth + Base
	sha256 := prime256v1.Sha256([]byte(salt))
	pwd := shift.Shift(salt, HashCode(string(sha256)))
	if len(pwd) > BaseKeyLen {
		pwd = pwd[:BaseKeyLen]
	}
	return pwd
}

// HashCode 将任何长度的字符串，通过运算，散列成0-15整数
func HashCode(key string) int {
	var index int = 0
	index = int(key[0])
	for k := 0; k < len(key); k++ {
		// 1103515245是个好数字，使通过hashCode散列出的0-15的数字的概率是相等的
		index *= 1103515245 + int(key[k])
	}
	index >>= 27
	index &= 16 - 1
	return index
}
