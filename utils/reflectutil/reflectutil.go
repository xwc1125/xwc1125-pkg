// Package reflectutil
//
// @author: xwc1125
package reflectutil

import (
	"fmt"
	"reflect"
	"strings"
)

const tagName = "tag"

const (
	extends    = "extends" // 对结构体使用，会延伸到对应的结构体中，将结构体中的字段赋值到最外层
	useTagName = "utn"     // 对结构体使用，保留原始的结构嵌套
)

func GetMapByTag(v reflect.Value, res map[string]interface{}, tagName string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("getMapByTag err: ", err)
		}
	}()

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}

		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		if structField.PkgPath != "" {
			// 小写的不导出
			// reflect.Value.Interface: cannot return value obtained from unexported field or method
			continue
		}
		valueField := v.Field(i)

		tag, ok := structField.Tag.Lookup(tagName)
		if !ok {
			continue
		}

		var mapKey string
		var mapV interface{}
		var extend bool // 结构体嵌套，放入最外层的map中
		tagSlice := strings.Split(tag, ",")
		for _, t := range tagSlice {
			if t == extends {
				GetMapByTag(valueField, res, tagName)
				extend = true
				break
			} else if t == useTagName {
				// 声明嵌套的map
				embeddedMap := make(map[string]interface{})
				GetMapByTag(valueField, embeddedMap, tagName)
				mapV = embeddedMap
			} else {
				mapKey = t
			}
		}
		if extend {
			continue
		}
		if mapKey == "" {
			continue
		}
		if mapV != nil {
			res[mapKey] = mapV
			continue
		}

		// 这里只做了结构体切片，其他的不管，直接放入mapV
		if valueField.Kind() == reflect.Slice {
			if valueField.Len() != 0 {
				var sliceDataValue = valueField.Index(0) // 取第一个，拿它的 type
				if sliceDataValue.Kind() == reflect.Ptr && !sliceDataValue.IsNil() {
					sliceDataValue = sliceDataValue.Elem()
				}
				if sliceDataValue.Kind() == reflect.Struct {
					newSlice := make([]map[string]interface{}, valueField.Len())
					for i := 0; i < valueField.Len(); i++ {
						oneMap := make(map[string]interface{})
						GetMapByTag(valueField.Index(i), oneMap, tagName)
						newSlice[i] = oneMap
					}
					mapV = newSlice
				} else {
					mapV = valueField.Interface()
				}
			} else {
				mapV = valueField.Interface()
			}
		} else {
			mapV = valueField.Interface()
		}

		res[mapKey] = mapV
	}
}
