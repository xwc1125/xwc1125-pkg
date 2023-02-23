// Package db_gorm
//
// @author: xwc1125
package db_gorm

import (
	"sync"
	"time"

	"github.com/xwc1125/xwc1125-pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var (
	masterEngine *gorm.DB
	lock         sync.Mutex
)

// MasterEngine 主库，单例
func MasterEngine(config database.MysqlConfig) *gorm.DB {
	if masterEngine != nil {
		return masterEngine
	}

	lock.Lock()
	defer lock.Unlock()

	if masterEngine != nil {
		return masterEngine
	}

	masterEngine = NewEngine(config)
	return masterEngine
}

func NewEngine(config database.MysqlConfig) *gorm.DB {
	driver, ok := GormDBOpens[config.Driver]
	if !ok {
		return nil
	}

	logLvl := logger.Error
	if config.ShowSQL {
		logLvl = logger.Info
	}

	strategy := schema.NamingStrategy{
		SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
	}
	if config.PrefixTable != "" {
		// 设置表前缀
		strategy.TablePrefix = config.PrefixTable
	}
	if config.PrefixColumn != "" {
		// 设置字段前缀
		strategy.NameReplacer = &replace{
			prefixColumn: config.PrefixColumn,
		}
	}

	db, err := gorm.Open(driver(config.DSN()), &gorm.Config{
		NamingStrategy: strategy,
		// 设置gorm的logger显示
		Logger: New(
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel:      logLvl,
			},
		),
	})
	if err != nil {
		return nil
	}
	var register *dbresolver.DBResolver
	if register == nil {
		register = dbresolver.Register(dbresolver.Config{
			// // `db2` 作为 sources，`db3`、`db4` 作为 replicas
			// Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
			// Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
			// // sources/replicas 负载均衡策略
			// Policy: dbresolver.RandomPolicy{},
		})
	}
	if config.ConnMaxIdleTime > 0 {
		register = register.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
	}
	if config.ConnMaxLifetime > 0 {
		register = register.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}
	if config.MaxOpenConns > 0 {
		register = register.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		register = register.SetMaxIdleConns(config.MaxIdleConns)
	}
	if register != nil {
		err = db.Use(register)
	}
	return db
}

// 中国时区
var SysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")

type replace struct {
	prefixColumn string
}

func (r *replace) Replace(name string) string {
	return r.prefixColumn + name
}
