// Package protocol
//
// @author: xwc1125
package protocol

const (
	KEY_REQUEST_OBJ = "KEY_REQUEST_OBJ"
)

var (
	DefaultProtocol = ProtocolConfig{
		IsProtocol:     true,
		FilterSignList: []string{"rsa", "sign"},
		WhiteApiList:   []string{}, // 不需要进行协议处理的api,
		Limit: LimitConfig{
			AntiBrushFlag: false,
			InterTime:     10,
		},
	}

	DefaultPrivateKey = `
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

	// 	publicKey = `-----BEGIN PUBLIC KEY-----
	// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDS/KWjwOjpryCztRlIRaIuCuO5
	// nRnYA/I+cRnKS9k6ciaaAFGUpTqC8CU8hJb6i1ggt6Pec21LB/BYorsMhFdq6UsO
	// 9436T772FN+E6DNUsI16bmHx07R2uwxG7QiKttMZLhFyMgu5pMP+qQAs6CZfzR5d
	// WdhILxUVVwA1GIKItQIDAQAB
	// -----END PUBLIC KEY-----`
)

type ProtocolConfig struct {
	IsProtocol     bool        `json:"is_protocol" mapstructure:"is_protocol"`
	FilterSignList []string    `json:"filter_sign_list" mapstructure:"filter_sign_list"`
	WhiteApiList   []string    `json:"white_api_list" mapstructure:"white_api_list"`
	Limit          LimitConfig `json:"limit" mapstructure:"limit"`
	PrivateKey     string      `json:"private_key" mapstructure:"private_key"`
}

type LimitConfig struct {
	AntiBrushFlag bool `json:"anti_brush_flag" mapstructure:"anti_brush_flag"`
	InterTime     int  `json:"inter_time" mapstructure:"inter_time"`
}
