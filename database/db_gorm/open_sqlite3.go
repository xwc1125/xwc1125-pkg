//go:build sqlite3
// +build sqlite3

package db_gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GormDBOpens = map[string]func(string) gorm.Dialector{
	"mysql":    mysql.Open,
	"postgres": postgres.Open,
	"sqlite3":  sqlite.Open,
}
