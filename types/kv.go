// Package types
//
// @author: xwc1125
package types

type KV[K any, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}
