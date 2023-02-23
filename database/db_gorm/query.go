// Package db_gorm
//
// @author: xwc1125
package db_gorm

import (
	"github.com/xwc1125/xwc1125-pkg/database/search"
	"gorm.io/gorm"
)

func MakeCondition(dbDrive string, q interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &search.OrmCondition{
			OrmWhere: search.OrmWhere{},
			Join:     make([]*search.OrmJoin, 0),
		}
		search.ResolveSearchQuery(dbDrive, q, condition)
		for _, join := range condition.Join {
			if join == nil {
				continue
			}
			db = db.Joins(join.JoinOn)
			for k, v := range join.Where {
				db = db.Where(k, v...)
			}
			for k, v := range join.Or {
				db = db.Or(k, v...)
			}
			for _, o := range join.Order {
				db = db.Order(o)
			}
		}
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		return db
	}
}
