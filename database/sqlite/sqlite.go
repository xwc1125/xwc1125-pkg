// Package sqlite
//
// @author: xwc1125
package sqlite

import (
	"database/sql"
	"time"

	"github.com/chain5j/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB ...
type DB struct {
	log logger.Logger
	db  *sqlx.DB
}

// New ...
func New(config *SqliteConfig) (*DB, error) {
	db, err := sqlx.Connect("mysql", config.Datasource)
	if err != nil {
		logger.Error("mysql connect err", "err", err)
		return nil, err
	}
	// 设置连接池最大连接数
	db.SetMaxOpenConns(100)
	// 设置连接池最大空闲连接数
	db.SetMaxIdleConns(20)
	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime))
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime))
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.ConnMaxIdleTime)
	}
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	return &DB{
		db:  db,
		log: logger.Log("Mysql"),
	}, nil
}

// DB ...
func (db *DB) DB() *sqlx.DB {
	return db.db
}

// Select 查询
func (db *DB) Select(dest interface{}, sql string, args ...interface{}) error {
	return db.db.Get(dest, sql, args...)
}

// Insert 插入
func (db *DB) Insert(sql string, args ...interface{}) (id int64, err error) {
	result, err := db.db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

// Update 更新
func (db *DB) Update(sql string, args ...interface{}) (affect int64, err error) {
	result, err := db.db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// Delete 删除
func (db *DB) Delete(sql string, args ...interface{}) (affect int64, err error) {
	result, err := db.db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// Transaction 事务处理
func (db *DB) Transaction(sqls []DBSql) (err error) {
	txs, err := db.db.Begin()
	if err != nil {
		db.log.Error("begin transaction err", "err", err)
		return err
	}
	// 如果出现了异常，导致没有 commit和rollback，可以用来收尾
	defer db.clearTransaction(txs)

	if sqls != nil {
		for _, f := range sqls {
			result, err := db.db.Exec(f.Sql, f.Args)
			if err != nil {
				db.log.Error("exec transaction err", "dbSql", f, "err", err)
				txs.Rollback()
				return err
			}
			_, err = result.RowsAffected()
			if err != nil {
				db.log.Error("affect err", "dbSql", f, "err", err)
				txs.Rollback()
				return err
			}
		}
	}

	if err != nil {
		txs.Rollback()
		db.log.Error("transaction err", "err", err)
		return
	}
	return txs.Commit()
}

// 事务回滚
func (db *DB) clearTransaction(txs *sql.Tx) error {
	db.log.Debug("clearTransaction rollback")
	err := txs.Rollback()
	if err != sql.ErrTxDone && err != nil {
		logger.Error("clearTransaction err", "err", err)
		return err
	}
	return nil
}
