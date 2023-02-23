// Package jwtauth
//
// @author: xwc1125
package jwtauth

import (
	"crypto/rsa"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

type JwtHandler struct {
	secret string

	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var (
	JwtHandlerPool = &sync.Pool{New: func() interface{} {
		return &JwtHandler{}
	}}
)

func NewJwtHandler() *JwtHandler {
	return JwtHandlerPool.Get().(*JwtHandler)
}

// GenerateToken 生成token
func (j *JwtHandler) GenerateToken(claims JwtClaims) (string, error) {
	var (
		token string
		err   error
	)
	if j.privateKey != nil {
		m := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token, err = m.SignedString(j.privateKey)
	} else {
		m := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err = m.SignedString([]byte(j.secret))
	}
	if err != nil {
		log().Error("generate token", "err", err)
	}
	return token, err
}

// GenerateRefreshToken 生成刷新token
func (j *JwtHandler) GenerateRefreshToken(claims JwtRefreshClaims) (string, error) {
	var (
		token string
		err   error
	)
	if j.privateKey != nil {
		m := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token, err = m.SignedString(j.privateKey)
	} else {
		m := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err = m.SignedString([]byte(j.secret))
	}
	if err != nil {
		log().Error("generate refresh token", "err", err)
	}
	return token, err

}

func (j *JwtHandler) SetPublicKey(key *rsa.PublicKey) {
	j.publicKey = key
}

func (j *JwtHandler) GetPublicKey() *rsa.PublicKey {
	return j.publicKey
}

func (j *JwtHandler) SetPrivateKey(key *rsa.PrivateKey) {
	j.privateKey = key
}

func (j *JwtHandler) GetPrivateKey() *rsa.PrivateKey {
	return j.privateKey
}
func (j *JwtHandler) SetSecret(secret string) {
	j.secret = secret
}

func (j *JwtHandler) GetSecret() string {
	return j.secret
}

// ValidateToken 验证token
func (j *JwtHandler) ValidateToken(token string) (*JwtClaims, error) {
	parsedClaims := &JwtClaims{}
	_, err := jwt.ParseWithClaims(token, parsedClaims, func(*jwt.Token) (interface{}, error) {
		if j.publicKey != nil {
			return j.publicKey, nil
		} else {
			return []byte(j.secret), nil
		}
	})
	if err != nil {
		log().Error("validate token", "err", err)
	}
	return parsedClaims, err
}

func (j *JwtHandler) ValidateRefreshToken(token string) (*JwtRefreshClaims, error) {
	parsedClaims := &JwtRefreshClaims{}
	_, err := jwt.ParseWithClaims(token, parsedClaims, func(*jwt.Token) (interface{}, error) {
		if j.publicKey != nil {
			return j.publicKey, nil
		} else {
			return []byte(j.secret), nil
		}
	})
	if err != nil {
		log().Error("validate refresh token", "err", err)
	}
	return parsedClaims, err
}

func (j *JwtHandler) Release() {
	JwtHandlerPool.Put(j)
}
