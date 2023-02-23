// Package db_xorm
//
// @author: xwc1125
package db_xorm

import (
	"reflect"

	"github.com/xwc1125/xwc1125-pkg/types/response"
	"xorm.io/xorm"
)

func SelectPage(db *xorm.Engine, page *response.PageInfo, wrapper map[string]interface{}, model interface{}) (e error) {
	e = nil
	db.Table(&model).Where(wrapper).Count(&page.Total)
	if page.Total == 0 {
		page.List = []interface{}{}
		return
	}
	// 反射获得类型
	t := reflect.TypeOf(model)
	// 再通过反射创建创建对应类型的数组
	list := reflect.Zero(reflect.SliceOf(t)).Interface()
	session := db.Table(&model).Where(wrapper)
	e = Paginate(page)(session).Find(&list)
	if e != nil {
		return
	}
	page.List = list
	return
}

func Paginate(page *response.PageInfo) func(db *xorm.Session) *xorm.Session {
	return func(db *xorm.Session) *xorm.Session {
		if page.PageIndex <= 0 {
			page.PageIndex = 1
		}
		switch {
		case page.PageSize > 100:
			page.PageSize = 100
		case page.PageSize <= 0:
			page.PageSize = 10
		}
		page.Pages = page.Total / page.PageSize
		if page.Total%page.PageSize != 0 {
			page.Pages++
		}
		p := page.PageIndex
		if page.PageIndex > page.Pages {
			p = page.Pages
		}
		size := page.PageSize
		offset := int((p - 1) * size)
		return db.Limit(int(size), offset)
	}
}
