// Package captcha
//
// @author: xwc1125
package captcha

import (
	"fmt"
	"testing"
)

func TestCaptcha(t *testing.T) {
	typs := []string{Unknown, String, Math, Chinese, Audio}
	for _, typ := range typs {
		fmt.Println("\n\n=================" + typ + "=================")
		id, code, err := GenerateDefault(typ)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("id=", id)
		fmt.Println("code=", code)
	}

}
