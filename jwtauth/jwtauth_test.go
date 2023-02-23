// Package jwtauth
//
// @author: xwc1125
package jwtauth

import (
	"testing"

	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap"
)

func init() {
	zap.InitWithConfig(&logger.LogConfig{
		Console: logger.ConsoleLogConfig{
			Level:    logger.LvlDebug,
			Modules:  "*",
			ShowPath: false,
			Format:   "",
			UseColor: true,
			Console:  true,
		},
		File: logger.FileLogConfig{},
	})
}

func TestGenerateToken(t *testing.T) {
	config := JWTConfig{
		Secret:           "xwc1125",
		Timeout:          3600,
		RefreshTimeout:   0,
		IgnoreURLs:       nil,
		AuthorizationKey: "Authorization",
		OnlineKey:        "",
		PriKeyPath:       "",
		PubKeyPath:       "",
	}
	token, err := GenerateToken(config, 1, "admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	jwtClaims, err := ParseToken(config, token)
	if err != nil {
		t.Fatal(err)
	}
	_ = jwtClaims

	refreshToken, err := GenerateRefreshToken(config, token)
	if err != nil {
		t.Fatal(err)
	}
	refreshClaims, err := ParseRefreshToken(config, refreshToken)
	if err != nil {
		t.Fatal(err)
	}
	_ = refreshClaims
}
