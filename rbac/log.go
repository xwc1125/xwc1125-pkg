// Package rbac
//
// @author: xwc1125
package rbac

import (
	"fmt"
	"strings"

	"github.com/chain5j/logger"
)

func log() logger.Logger {
	return logger.Log("rbac")
}

type Logger struct {
	enabled bool
}

func (l *Logger) EnableLog(enable bool) {
	l.enabled = enable
}

func (l *Logger) IsEnabled() bool {
	return l.enabled
}

func (l *Logger) LogModel(model [][]string) {
	if !l.enabled {
		return
	}
	var str strings.Builder
	str.WriteString("Model: ")
	for _, v := range model {
		str.WriteString(fmt.Sprintf("%v\n", v))
	}
	log().Info(str.String())
}

func (l *Logger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	if !l.enabled {
		return
	}

	var reqStr strings.Builder
	reqStr.WriteString("Request: ")
	for i, rval := range request {
		if i != len(request)-1 {
			reqStr.WriteString(fmt.Sprintf("%v, ", rval))
		} else {
			reqStr.WriteString(fmt.Sprintf("%v", rval))
		}
	}
	reqStr.WriteString(fmt.Sprintf(" ---> %t\n", result))

	reqStr.WriteString("Hit Policy: ")
	for i, pval := range explains {
		if i != len(explains)-1 {
			reqStr.WriteString(fmt.Sprintf("%v, ", pval))
		} else {
			reqStr.WriteString(fmt.Sprintf("%v \n", pval))
		}
	}

	log().Info(reqStr.String())
}

func (l *Logger) LogPolicy(policy map[string][][]string) {
	if !l.enabled {
		return
	}

	var str strings.Builder
	str.WriteString("Policy: \n")
	for k, v := range policy {
		str.WriteString(fmt.Sprintf("%s : %v\n", k, v))
	}
	log().Info(str.String())
}

func (l *Logger) LogRole(roles []string) {
	if !l.enabled {
		return
	}

	log().Info("Roles: ", roles)
}
