// Package db_gorm
//
// @author: xwc1125
package db_gorm

import (
	"reflect"

	"github.com/xwc1125/xwc1125-pkg/types"
	"github.com/xwc1125/xwc1125-pkg/types/response"
	"gorm.io/gorm"
)

func SelectPage(db *gorm.DB, req types.PageQuery, wrapper map[string]interface{}, model interface{}) (result *response.PageInfo, err error) {
	db.Model(&model).
		Where(wrapper).
		Count(&result.Total)
	if result.Total == 0 {
		result.List = []interface{}{}
		return
	}
	// 反射获得类型
	t := reflect.TypeOf(model)
	// 再通过反射创建创建对应类型的数组
	list := reflect.Zero(reflect.SliceOf(t)).Interface()
	err = db.Model(&model).
		Where(wrapper).
		Scopes(MakePaginate(req)).
		Find(&list).
		Error
	if err != nil {
		return
	}
	result.List = list
	result.PageIndex = req.GetPageIndex()
	result.PageSize = req.GetPageSize()
	result.Pages = result.Total / result.PageSize
	if result.Total%result.PageSize != 0 {
		result.Pages++
	}
	return
}

func MakePaginate(page types.PageQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		startTime := page.GetStartTime()
		endTime := page.GetEndTime()
		if endTime > 0 {
			db = db.Where("created_at <= ?", endTime)
		}
		if startTime > 0 {
			db = db.Where("created_at >= ?", startTime)
		}
		db = db.Order(page.GetOrder())
		db = db.Offset(int(page.GetOffset())).Limit(int(page.GetPageSize()))
		return db
	}
}
func Paginate(pageSize, pageIndex int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}
