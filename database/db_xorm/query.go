// Package db_xorm
//
// @author: xwc1125
package db_xorm

import (
	"github.com/xwc1125/xwc1125-pkg/database/search"
	"xorm.io/xorm"
)

func MakeCondition(dbDrive string, q interface{}) func(db *xorm.Session) *xorm.Session {
	return func(db *xorm.Session) *xorm.Session {
		condition := &search.OrmCondition{
			OrmWhere: search.OrmWhere{},
			Join:     make([]*search.OrmJoin, 0),
		}
		search.ResolveSearchQuery(dbDrive, q, condition)
		for _, join := range condition.Join {
			if join == nil {
				continue
			}
			db = db.Join("on", join.JoinOn, "")
			for k, v := range join.Where {
				db = db.Where(k, v...)
			}
			for k, v := range join.Or {
				db = db.Or(k, v...)
			}
			for _, o := range join.Order {
				db = db.OrderBy(o)
			}
		}
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for _, o := range condition.Order {
			db = db.OrderBy(o)
		}
		return db
	}
}
