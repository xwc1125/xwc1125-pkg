// Package ldap
//
// @author: xwc1125
package ldap

import "fmt"

var (
	internalServerError = fmt.Errorf("internal server error, try again later please")
	loginFailError      = fmt.Errorf("login fail, check your username and password")
)
