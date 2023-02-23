// Package db_gorm
//
// @author: xwc1125
package db_gorm

import (
	"context"
	"fmt"
	"time"

	loggerCore "github.com/chain5j/logger"
	"github.com/xwc1125/xwc1125-pkg/protocol/contextx"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	_ logger.Interface = new(gormLogger)
)

type gormLogger struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func New(config logger.Config) logger.Interface {
	return &gormLogger{
		Config: config,
	}
}

func (l *gormLogger) getLogger(ctx context.Context) loggerCore.Logger {
	requestId := ctx.Value(contextx.RequestIDKey)
	if requestId != nil {
		return loggerCore.New("gorm", "x-request-id", requestId)
	}
	return loggerCore.New("gorm")
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log := l.getLogger(ctx)
		log.Info(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log := l.getLogger(ctx)
		log.Warn(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log := l.getLogger(ctx)
		log.Error(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > logger.Silent {
		log := l.getLogger(ctx)
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			sql, rows := fc()
			log.Trace(sql, "rows", rows, "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6)
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			log.Trace(sql, "rows", rows, "slowLog", slowLog, "elapsed", float64(elapsed.Nanoseconds())/1e6)
		case l.LogLevel == logger.Info:
			sql, rows := fc()
			log.Trace(sql, "rows", rows, "elapsed", float64(elapsed.Nanoseconds())/1e6)
		}
	}
}

type traceRecorder struct {
	logger.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

func (l traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface, BeginAt: time.Now()}
}

func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}
