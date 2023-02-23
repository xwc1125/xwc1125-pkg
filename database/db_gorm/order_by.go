// Package db_gorm
//
// @author: xwc1125
package db_gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func OrderDest(sort string, bl bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: sort}, Desc: bl})
	}
}
