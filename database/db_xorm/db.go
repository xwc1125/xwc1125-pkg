package db_xorm

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xwc1125/xwc1125-pkg/database"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	masterEngine *xorm.Engine
	slaveEngine  *xorm.Engine
	lock         sync.Mutex
)

// MasterEngine 主库，单例
func MasterEngine(config database.MysqlConfig) *xorm.Engine {
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

func NewEngine(config database.MysqlConfig) *xorm.Engine {
	engine, err := xorm.NewEngine(config.Driver, config.DSN())
	if err != nil {
		log().Error("new db engine error!!", "err", err)
		return nil
	}
	settings(engine, config)
	return engine
}

// 中国时区
var SysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")

func settings(engine *xorm.Engine, config database.MysqlConfig) {
	engine.ShowSQL(config.ShowSQL)
	engine.SetTZLocation(SysTimeLocation)
	if config.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.ConnMaxLifetime > 0 {
		engine.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}
	if config.PrefixTable != "" {
		// 设置表前缀
		// engine.SetTableMapper(core.SameMapper{})
		// engine.SetColumnMapper(core.SnakeMapper{})
		tbMapper := core.NewPrefixMapper(core.SameMapper{}, config.PrefixTable)
		engine.SetTableMapper(tbMapper)
	} else {
		engine.SetTableMapper(core.GonicMapper{})
	}
	if config.PrefixColumn != "" {
		// 设置字段前缀
		tbMapper := core.NewPrefixMapper(core.SameMapper{}, config.PrefixColumn)
		engine.SetColumnMapper(tbMapper)
	} else {
		engine.SetTableMapper(core.GonicMapper{})
	}

	// 性能优化的时候才考虑，加上本机的SQL缓存
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// engine.SetDefaultCacher(cacher)
}
