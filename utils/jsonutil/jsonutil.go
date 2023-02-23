// Package jsonutil
//
// @author: xwc1125
package jsonutil

import (
	"bytes"
	"io"
	"net/http"

	"github.com/chain5j/chain5j-pkg/codec/json"
)

func ParsePayload(c *http.Request, obj interface{}) error {
	buf := make([]byte, 1024)
	n, _ := c.Body.Read(buf)
	// s := string(buf[0:n])
	return BindBody(buf[0:n], obj)
}

func BindBody(body []byte, obj interface{}) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

var EnableDecoderUseNumber = false

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}

	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
