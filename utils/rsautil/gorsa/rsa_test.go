package gorsa

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func init() {

}

var Pubkey = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDS/KWjwOjpryCztRlIRaIuCuO5
nRnYA/I+cRnKS9k6ciaaAFGUpTqC8CU8hJb6i1ggt6Pec21LB/BYorsMhFdq6UsO
9436T772FN+E6DNUsI16bmHx07R2uwxG7QiKttMZLhFyMgu5pMP+qQAs6CZfzR5d
WdhILxUVVwA1GIKItQIDAQAB
-----END PUBLIC KEY-----
`

var Pirvatekey = `
-----BEGIN RSA PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBANL8paPA6OmvILO1
GUhFoi4K47mdGdgD8j5xGcpL2TpyJpoAUZSlOoLwJTyElvqLWCC3o95zbUsH8Fii
uwyEV2rpSw73jfpPvvYU34ToM1SwjXpuYfHTtHa7DEbtCIq20xkuEXIyC7mkw/6p
ACzoJl/NHl1Z2EgvFRVXADUYgoi1AgMBAAECgYA36r+Xg7P24vQezJtTvSi7XSV3
dvx3dKxKAj2ckLeUttYmMHlulM/KDg1AWb+NzUpg+WqTtTG9FfCL/Sznp1dnQah0
XS3vsuh20V3xKKBVoh/cI63/qRgcDBSPn7t9qNZ2SX0Rl1SnGxiE7TcJK/f3H/uK
IeFQgqG0hYvfSiScYQJBAPyrTMhrQjEo3QW5p2L18UpHjhEU5IYS+Ax2IyIcrAtL
uQ9JFZwFW+yXAJNDTfzBdS4jI7BIXsA85GRrZU0GNw0CQQDVxK3dVI2hJ+Yu1xA6
iSg7FLTKdIwV/q/VWwbEdsr9HyAYXbHu4dJpOZgNKQisvLM0GYqd8Pn8Tii7+4tV
Gq5JAkEAs857j8y0iCNaVm6t7cCz+3Y8ZW+GyNrK5qNTkTzyOf+jLHuIA0XVCuLS
p/mnkA1kBHdBOHvn4cnzhnre1hdsKQJAdzjNGx7YKqQ9DaymgW8Tf/fpaOydYHr+
CAlPee0jAw8D8HL5FNjfaA5WDijvjJ9lds4z8CiA08WnlEgTinBp+QJBAKEReXrM
jG/mqNFa/G8lVYm/y5NuZrswshleL5y9aw6wfyik75V3cBTZVmW6BgnRiLaJiwpu
DiNQ7YgWwzk+zDY=
-----END RSA PRIVATE KEY-----
`

func TestRSA(t *testing.T) {
	// 公钥加密私钥解密
	if err := applyPubEPriD(); err != nil {
		log.Println(err)
	}
	// 公钥解密私钥加密
	if err := applyPriEPubD(); err != nil {
		log.Println(err)
	}
}

// 公钥加密私钥解密
func applyPubEPriD() error {
	pubenctypt, err := PublicEncrypt(`hello world`, Pubkey)
	if err != nil {
		return err
	}

	pridecrypt, err := PriKeyDecrypt(pubenctypt, Pirvatekey)
	if err != nil {
		return err
	}
	fmt.Println(string(pridecrypt))
	if string(pridecrypt) != `hello world` {
		return errors.New(`解密失败`)
	}
	return nil
}

// 公钥解密私钥加密
func applyPriEPubD() error {
	prienctypt, err := PriKeyEncrypt(`hello world`, Pirvatekey)
	if err != nil {
		return err
	}

	pubdecrypt, err := PublicDecrypt(prienctypt, Pubkey)
	if err != nil {
		return err
	}
	if string(pubdecrypt) != `hello world` {
		return errors.New(`解密失败`)
	}
	return nil
}
