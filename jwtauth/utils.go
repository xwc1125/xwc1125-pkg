// Package jwtauth
//
// @author: xwc1125
package jwtauth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chain5j/chain5j-pkg/codec/json"
	"github.com/dgrijalva/jwt-go"
)

var (
	keyMap sync.Map
)

// GenerateToken 生成token
func GenerateToken(config JWTConfig, userId int64, userName string, extraData MapClaims) (string, error) {
	jwtHandler := getJwtHandler(config)
	defer jwtHandler.Release()

	expireTime := time.Now().Add(time.Duration(config.Timeout) * time.Second)
	if extraData == nil {
		extraData = MapClaims{}
	}
	extraData["exp"] = expireTime.Unix()

	jwtClaims := JwtClaims{
		UserToken{
			Uid:       userId,
			Username:  userName,
			ExtraData: extraData,
		},
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    Issuer,
		},
	}
	return jwtHandler.GenerateToken(jwtClaims)
}

func getJwtHandler(config JWTConfig) *JwtHandler {
	jwtHandler := NewJwtHandler()
	jwtHandler.SetSecret(config.Secret)
	if len(config.PriKeyPath) > 0 && len(config.PubKeyPath) > 0 {
		jwtPrivateKey, err1 := filepath.Abs(config.PriKeyPath)
		privateKey, err2 := LoadRSAPrivateKey(jwtPrivateKey)
		jwtPubKeyKey, err3 := filepath.Abs(config.PubKeyPath)
		publicKey, err4 := LoadRSAPublicKey(jwtPubKeyKey)

		if err1 == nil && err2 == nil || err3 == nil && err4 == nil {
			jwtHandler.SetPrivateKey(privateKey)
			jwtHandler.SetPublicKey(publicKey)
		}
	}
	return jwtHandler
}

// GenerateRefreshToken 生成刷新token
func GenerateRefreshToken(config JWTConfig, token string) (string, error) {
	jwtHandler := getJwtHandler(config)
	defer jwtHandler.Release()

	expireTime := time.Now().Add(time.Duration(config.Timeout) * time.Second)
	jwtClaims := JwtRefreshClaims{
		Token: token,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    Issuer,
		},
	}
	return jwtHandler.GenerateRefreshToken(jwtClaims)
}

// ParseToken 解析token
func ParseToken(config JWTConfig, token string) (*JwtClaims, error) {
	jwtHandler := getJwtHandler(config)
	// 格式化token
	token = TokenFormat(token)
	if len(token) == 0 {
		return nil, fmt.Errorf(TokenParseFailedAndEmpty)
	}
	return jwtHandler.ValidateToken(token)
}

// ParseRefreshToken 解析refresh token
func ParseRefreshToken(config JWTConfig, token string) (*JwtRefreshClaims, error) {
	jwtHandler := getJwtHandler(config)
	// 格式化token
	token = TokenFormat(token)
	if len(token) == 0 {
		return nil, fmt.Errorf(TokenParseFailedAndEmpty)
	}
	return jwtHandler.ValidateRefreshToken(token)
}

// TokenFormat 格式化token
func TokenFormat(token string) string {
	format, err := tokenFormat(token)
	if err != nil {
		log().Error("format token err", "err", err)
	}
	return format
}

func tokenFormat(token string) (string, error) {
	if token == "" {
		return "", nil
	}
	// TODO: Make this a bit more robust, parsing-wise
	token = strings.TrimSpace(token)
	token = strings.ReplaceAll(token, "\r", "")
	token = strings.ReplaceAll(token, "\n", "")
	authHeaderParts := strings.Split(token, " ")
	if len(authHeaderParts) == 1 {
		return authHeaderParts[0], nil
	}
	if len(authHeaderParts) == 2 {
		if strings.Trim(authHeaderParts[0], " ") == "" || strings.Trim(authHeaderParts[0], " ") == "Bearer" {
			return authHeaderParts[1], nil
		}
	}
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// MapClaimsToJwt map转jwtClaims
func MapClaimsToJwt(claims jwt.MapClaims) *JwtClaims {
	var jwtClaims *JwtClaims
	bytes, err := json.Marshal(claims)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &jwtClaims)
	if err != nil {
		return nil
	}
	return jwtClaims
}

// TokenToMapClaims jwt token 转mapClaims
func TokenToMapClaims(token *jwt.Token) jwt.MapClaims {
	if token == nil {
		return make(jwt.MapClaims)
	}

	claims := jwt.MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims
}

// LoadRSAPrivateKey 加载rsa私钥
func LoadRSAPrivateKey(prvFile string) (*rsa.PrivateKey, error) {
	if key, set := keyMap.Load(prvFile); set {
		return key.(*rsa.PrivateKey), nil
	} else {
		keyData, e := ioutil.ReadFile(prvFile)
		if e != nil {
			log().Error("load rsa privateKey from disk", "err", e)
			return nil, e
		}
		key, e := jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if e != nil {
			log().Error("parse rsa privateKey from pem", "err", e)
			return nil, e
		}
		keyMap.Store(prvFile, key)
		return key, nil
	}
}

// LoadRSAPublicKey 加载rsa公钥
func LoadRSAPublicKey(pubFile string) (*rsa.PublicKey, error) {
	if key, set := keyMap.Load(pubFile); set {
		return key.(*rsa.PublicKey), nil
	} else {
		keyData, e := ioutil.ReadFile(pubFile)
		if e != nil {
			log().Error("load rsa publicKey from disk", "err", e)
			return nil, e
		}
		key, e := jwt.ParseRSAPublicKeyFromPEM(keyData)
		if e != nil {
			log().Error("parse rsa publicKey from pem", "err", e)
			return nil, e
		}
		keyMap.Store(pubFile, key)
		return key, nil
	}
}

// FromAuthHeader 从Authorization header中获取
func FromAuthHeader(authKey string) TokenExtractor {
	return func(req *http.Request) (string, error) {
		authHeader := req.Header.Get(authKey)
		return tokenFormat(authHeader)
	}
}

// FromParameter 从params中获取
func FromParameter(tokenKey string) TokenExtractor {
	return func(req *http.Request) (string, error) {
		token := req.URL.Query().Get(tokenKey)
		return tokenFormat(token)
	}
}

// FromCookie 从cookie中获取
func FromCookie(param string) TokenExtractor {
	return func(req *http.Request) (string, error) {
		cookie, err := req.Cookie(param)
		if err != nil {
			return "", nil
		}
		token, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			return "", err
		}
		return tokenFormat(token)
	}
}

// FromFirst 获取它找到的第一个令牌
func FromFirst(extractors ...TokenExtractor) TokenExtractor {
	return func(req *http.Request) (string, error) {
		for _, ex := range extractors {
			token, err := ex(req)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", nil
	}
}
