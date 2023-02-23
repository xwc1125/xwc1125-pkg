// Package jsonp
// Package
//
// @author: xwc1125
package jsonp

import (
	"fmt"
	"net/http"

	"github.com/chain5j/chain5j-pkg/codec/json"
)

var (
	JsonpKey = "callback"
)

// IsJsonp 判断请求是否为jsonp
func IsJsonp(r *http.Request) bool {
	callbackName := r.URL.Query().Get(JsonpKey)
	if callbackName == "" {
		return false
	}
	return true
}

// JsonpHandler jsonp处理
func JsonpHandler(w http.ResponseWriter, r *http.Request) {
	callbackName := r.URL.Query().Get(JsonpKey)
	if callbackName == "" {
		fmt.Fprintf(w, "Please give callback name in query string")
		return
	}

	b, err := json.Marshal(r.Header)
	if err != nil {
		fmt.Fprintf(w, "json encode error")
		return
	}

	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(w, "%s(%s);", callbackName, b)
}
