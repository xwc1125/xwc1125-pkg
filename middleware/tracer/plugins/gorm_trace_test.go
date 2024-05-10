// Package plugins
package plugins

// import (
// 	"context"
// 	"fmt"
// 	"testing"
//
// 	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
// 	"github.com/uptrace/opentelemetry-go-extra/otelplay"
// 	"github.com/xwc1125/xwc1125-pkg/database"
// 	"github.com/xwc1125/xwc1125-pkg/database/db_gorm"
// 	"go.opentelemetry.io/otel"
// )
//
// type Test1 struct {
// 	database.ModelID
// 	database.ModelTime `xorm:"extends" gorm:"extends"`
// }
//
// func TestNewTracerForGorm(t *testing.T) {
// 	ctx := context.Background()
//
// 	shutdown := otelplay.ConfigureOpentelemetry(ctx)
// 	defer shutdown()
//
// 	engine := db_gorm.NewEngine(database.MysqlConfig{
// 		Driver:          "mysql",
// 		Url:             "tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local&timeout=1000ms",
// 		Username:        "root",
// 		Password:        "root_@123456",
// 		Secret:          "",
// 		MaxIdleConns:    1000,
// 		MaxOpenConns:    1000,
// 		ConnMaxLifetime: 100,
// 		ConnMaxIdleTime: 100,
// 		PrefixTable:     "",
// 		PrefixColumn:    "",
// 		ShowSQL:         true,
// 		LogLevel:        3,
// 	})
// 	if err := engine.Use(otelgorm.NewPlugin()); err != nil {
// 		panic(err)
// 	}
// 	engine.Migrator().CreateTable(Test1{})
// 	var test1 Test1
// 	tx := engine.WithContext(ctx).Select("id", "1").Find(&test1)
// 	// 更新时，使用这个检测
// 	// if tx.RowsAffected == 0 {
// 	// 	t.Fatal("无影响")
// 	// }
// 	if tx.Error != nil {
// 		t.Fatal(tx.Error)
// 	}
// 	fmt.Println(test1)
//
// 	tracer := otel.Tracer("app_or_package_name")
//
// 	ctx, span := tracer.Start(ctx, "root")
// 	defer span.End()
//
// 	if err := engine.WithContext(ctx).Select("id", "1").Find(&test1).Error; err != nil {
// 		panic(err)
// 	}
//
// 	otelplay.PrintTraceID(ctx)
// }
