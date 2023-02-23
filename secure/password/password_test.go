// Package password
//
// @author: xwc1125
package password

import (
	"encoding/base64"
	"testing"
)

func TestPassword(t *testing.T) {
	hash, err := GeneratePassword("123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(hash))
	ok := CheckPassword("$2a$10$zOZ2Tl6qlCTceD.Sva868e2DrDzCjcwaoN7ZEtjLrvi4aiXdZflRa", "123456")
	t.Log(ok)
}

func TestEndPassword(t *testing.T) {
	var auth = "xwc1125"
	encKey, err := EncryptKey(Key{
		Id:         "1",
		PrivateKey: []byte("UxSZY*8462xv"),
	}, auth)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(encKey))
	pwdBase64 := base64.StdEncoding.EncodeToString(encKey)
	t.Log(string(pwdBase64))
	decryptKey, err := DecryptKey(encKey, auth)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(decryptKey.PrivateKey))
}
