// Package kvutil
//
// @author: xwc1125
package kvutil

import (
	"bytes"
	"net/url"
	"sort"
)

type kv struct {
	k string
	v string
}

type KeyValue struct {
	u url.Values
	s []kv
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		u: make(map[string][]string),
		s: make([]kv, 0, 2),
	}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (s *KeyValue) Get(key string) string {
	vs := s.u[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
func (s *KeyValue) Set(key, value string) *KeyValue {
	s.u[key] = []string{value}
	s.s = append(s.s, kv{k: key, v: value})
	return s
}

// Encode 参数编号返回编码后的字符串
func (s *KeyValue) Encode() string {
	return s.u.Encode()
}

// Sort 对参数进行参数
func (s *KeyValue) Sort() *KeyValue {
	sort.Slice(s.s, func(i, j int) bool {
		return s.s[i].k < s.s[j].k
	})
	return s
}

// JoinAll 对参数进行拼接
func (s *KeyValue) JoinAll(a string, b string) string {
	buffer := bytes.NewBufferString("")
	for i, v := range s.s {
		buffer.WriteString(v.k)
		buffer.WriteString(a)
		buffer.WriteString(v.v)
		if i < len(s.s)-1 {
			buffer.WriteString(b)
		}
	}
	return buffer.String()
}

// Join 只拼接值不为空的参数
func (s *KeyValue) Join(a string, b string) string {
	buffer := bytes.NewBufferString("")
	for i, v := range s.s {
		if v.v != "" {
			buffer.WriteString(v.k)
			buffer.WriteString(a)
			buffer.WriteString(v.v)
			if i < len(s.s)-1 {
				buffer.WriteString(b)
			}
		}

	}
	return buffer.String()
}

func (s *KeyValue) JoinValues() string {
	buffer := bytes.NewBufferString("")
	for _, v := range s.s {
		buffer.WriteString(v.v)
	}
	return buffer.String()
}
