// Package pmodel
//
// @author: xwc1125
package pmodel

import "github.com/chain5j/chain5j-pkg/codec/json"

type CoreDataInfo map[string]string

func (c *CoreDataInfo) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}
