// Package jwtauth
//
// @author: xwc1125
package jwtauth

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/xwc1125/xwc1125-pkg/protocol/contextx"
	"github.com/xwc1125/xwc1125-pkg/types/response"
)

type (
	// TokenExtractor is a function that takes a context as input and returns
	// either a token or an error.  An error should only be returned if an attempt
	// to specify a token was found, but the information was somehow incorrectly
	// formed.  In the case where a token is simply not present, this should not
	// be treated as an error.  An empty string should be returned in that case.
	TokenExtractor func(req *http.Request) (string, error)

	MapClaims jwt.MapClaims
)

var (
	jwtAuth *JwtAuth
	lock    sync.Mutex

	Issuer = "key-casbins-jwt"
)

// JwtAuth jwt auth
type JwtAuth struct {
	JWTConfig JWTConfig
	*JwtHandler
	Extractor TokenExtractor // Extractor 从请求中提取令牌的函数
}

// NewJwtAuth jwt中间件配置
func NewJwtAuth(config JWTConfig) *JwtAuth {
	if jwtAuth != nil {
		return jwtAuth
	}

	lock.Lock()
	defer lock.Unlock()

	if jwtAuth != nil {
		return jwtAuth
	}
	jwtHandler := getJwtHandler(config)
	defer jwtHandler.Release()
	tokenExtractors := make([]TokenExtractor, 0)
	if len(config.AuthorizationKey) == 0 {
		tokenExtractors = append(tokenExtractors, FromAuthHeader(AuthorizationKEY))
	} else {
		tokenExtractors = append(tokenExtractors, FromAuthHeader(config.AuthorizationKey))
	}
	if len(config.ParamTokenKey) > 0 {
		tokenExtractors = append(tokenExtractors, FromParameter(config.ParamTokenKey))
	}
	if len(config.CookieTokenKey) > 0 {
		tokenExtractors = append(tokenExtractors, FromCookie(config.CookieTokenKey))
	}

	jwtAuth = &JwtAuth{
		JWTConfig:  config,
		JwtHandler: jwtHandler,
		Extractor:  FromFirst(tokenExtractors...),
	}
	return jwtAuth
}

// GenerateToken 在登录成功的时候生成token
func (m *JwtAuth) GenerateToken(userId int64, userName string, extraData MapClaims) (string, error) {
	return GenerateToken(m.JWTConfig, userId, userName, extraData)
}

// CheckJWT 验证jwt
func (m *JwtAuth) CheckJWT(ctx contextx.Context) (*UserToken, error) {
	if !m.JWTConfig.EnableAuthOnOptions {
		if ctx.Request().Method == response.MethodOptions {
			return nil, TokenUnsupportedOptions
		}
	}

	// 获取token
	var (
		token string
		err   error
	)

	token, err = m.Extractor(ctx.Request())
	if err != nil {
		m.logf("error extracting JWT: %v", err)
		return nil, fmt.Errorf("error extracting token: %v", err)
	}

	if token == "" {
		m.logf("Error: No credentials found")
		return nil, fmt.Errorf(TokenParseFailedAndEmpty)
	}

	jwtClaims, err := ParseToken(m.JWTConfig, token)
	if err != nil {
		return nil, err
	}
	ctx.Set(m.JWTConfig.JwtContextKey+"_raw", token)
	return &jwtClaims.UserToken, nil
}

// logf 打印日志
func (m *JwtAuth) logf(format string, args ...interface{}) {
	if m.JWTConfig.Debug {
		log().Debug(fmt.Sprintf(format, args...))
	}
}

func Serve(ctx contextx.Context, config JWTConfig) (userToken *UserToken, b bool) {
	NewJwtAuth(config)
	var err error
	if userToken, err = jwtAuth.CheckJWT(ctx); err != nil {
		log().Error("Check jwt error", "err", err)
		return userToken, false
	}
	return userToken, true
}

func CheckPermission(ctx contextx.Context, config JWTConfig) bool {
	NewJwtAuth(config)
	if _, err := jwtAuth.CheckJWT(ctx); err != nil {
		log().Error("Check jwt error", "err", err)
		return false
	}
	return true
}

func GetUserId(config JWTConfig, token string) (int64, bool) {
	userToken, _ := ParseToken(config, token)
	if userToken != nil {
		return userToken.Uid, true
	}
	return -1, false
}
