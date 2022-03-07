// Package data
// @author: xs
// @date: 2022/3/7
// @Description: data db 重写日志数据库,使用zap 其他依赖有需要可以仿写
package data

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

const msgKey = "sql"
const msgInfo = "sql-info"

var (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

type log struct {
	zapLog        *zap.Logger
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration // 慢查询阀值
}

// 依赖
func NewDBLog(l *zap.Logger) *log {
	return &log{
		zapLog:        l,
		LogLevel:      logger.Info,     //打印所有日志，登记越小打印越少
		SlowThreshold: 2 * time.Second, // 2秒
	}
}

// LogMode log mode
func (l *log) LogMode(level logger.LogLevel) logger.Interface {
	newLog := *l
	newLog.LogLevel = level
	return &newLog
}

// Info print info
func (l log) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.zapLog.Info(msgKey,
			zap.String(msgInfo, fmt.Sprintf(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
		)
	}
}

// Warn print warn messages
func (l log) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.zapLog.Warn(msgKey,
			zap.String(msgInfo, fmt.Sprintf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
		)
	}
}

// Error print error messages
func (l log) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.zapLog.Error(msgKey, zap.String(msgInfo,
			fmt.Sprintf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...),
		))
	}
}

// Trace print sql message
func (l log) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	//|| !l.IgnoreRecordNotFoundError
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.zapLog.Error(msgKey, zap.String(msgInfo, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)))
		} else {
			l.zapLog.Error(msgKey, zap.String(msgInfo, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)))
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.zapLog.Warn(msgKey, zap.String(msgInfo, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)))
		} else {
			l.zapLog.Warn(msgKey, zap.String(msgInfo, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)))
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		fmt.Printf("%v", sql)
		if rows == -1 {
			l.zapLog.Info(msgKey, zap.String(msgInfo, fmt.Sprintf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)))
		} else {
			l.zapLog.Info(msgKey, zap.String(msgInfo, fmt.Sprintf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)))
		}
	}
}
