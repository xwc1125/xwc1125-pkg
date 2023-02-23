// Package captcha
//
// @author: xwc1125
package captcha

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

const (
	Unknown = "unknown"
	Audio   = "audio"
	String  = "string"
	Math    = "math"
	Chinese = "chinese"
)

// Driver 驱动
type Driver struct {
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var (
	Height     = 80
	Width      = 240
	Length     = 6
	Color      = &color.RGBA{R: 254, G: 254, B: 254, A: 254}
	NoiseCount = 80

	DefaultMemStore    = base64Captcha.DefaultMemStore
	DefaultDriverAudio = base64Captcha.DefaultDriverAudio
)

var (
	DefaultDriver = &Driver{
		DriverAudio: DefaultDriverAudio,
		DriverString: base64Captcha.NewDriverString(Height, Width, NoiseCount, base64Captcha.OptionShowHollowLine,
			Length, base64Captcha.TxtSimpleCharaters, Color, nil, nil),
		DriverChinese: base64Captcha.NewDriverChinese(Height, Width, NoiseCount, base64Captcha.OptionShowHollowLine,
			Length, base64Captcha.TxtChineseCharaters, Color, nil, nil),
		DriverMath: base64Captcha.NewDriverMath(Height, Width, NoiseCount, base64Captcha.OptionShowHollowLine,
			Color, nil, nil),
		DriverDigit: base64Captcha.NewDriverDigit(Height, Width, Length, 0.7, 80),
	}
)

// Generate 生成验证码
func Generate(store Store, captchaType string) (id, b64s string, err error) {
	var driver base64Captcha.Driver
	switch captchaType {
	case Audio:
		driver = DefaultDriver.DriverAudio
	case String:
		driver = DefaultDriver.DriverString.ConvertFonts()
	case Math:
		driver = DefaultDriver.DriverMath.ConvertFonts()
	case Chinese:
		driver = DefaultDriver.DriverChinese.ConvertFonts()
	default:
		driver = DefaultDriver.DriverDigit
	}
	c := base64Captcha.NewCaptcha(driver, store)
	return c.Generate()
}

func GenerateDefault(captchaType string) (id, b64s string, err error) {
	return Generate(DefaultMemStore, captchaType)
}

// Verify 验证验证码
func Verify(store Store, id, answer string, clear bool) bool {
	return store.Verify(id, answer, clear)
}

func VerifyDefault(id, answer string, clear bool) bool {
	return DefaultMemStore.Verify(id, answer, clear)
}
