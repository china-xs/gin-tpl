// Package data
// @author: xs
// @date: 2022/3/7
// @Description: data db 重写日志数据库,使用zap 其他依赖有需要可以仿写
package data

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

const msgKey = "sql"
const msgInfo = "sql-info"
const msgTrace = "trace_id"

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
		trace := zap.String(msgTrace, getTrace(ctx))
		l.zapLog.Info(msgKey,
			zap.String(msgInfo, fmt.Sprintf(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
			trace,
		)
	}
}

// Warn print warn messages
func (l log) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		trace := zap.String(msgTrace, getTrace(ctx))
		l.zapLog.Warn(msgKey,
			zap.String(msgInfo, fmt.Sprintf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
			trace,
		)
	}
}

// Error print error messages
func (l log) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		trace := zap.String(msgTrace, getTrace(ctx))
		l.zapLog.Error(msgKey, zap.String(msgInfo,
			fmt.Sprintf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
			trace)
	}
}

// Trace print sql message
func (l log) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	trace := zap.String(msgTrace, getTrace(ctx))

	elapsed := time.Since(begin)
	switch {
	//|| !l.IgnoreRecordNotFoundError
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.zapLog.Error(
				msgKey,
				zap.String(msgInfo, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)),
				trace,
			)
		} else {
			l.zapLog.Error(
				msgKey,
				zap.String(msgInfo, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)),
				trace,
			)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.zapLog.Warn(
				msgKey,
				zap.String(msgInfo, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)),
				trace,
			)
		} else {
			l.zapLog.Warn(
				msgKey,
				zap.String(msgInfo, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)),
				trace,
			)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()

		if rows == -1 {
			l.zapLog.Info(
				msgKey,
				zap.String(msgInfo, fmt.Sprintf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)),
				trace,
			)
		} else {
			l.zapLog.Info(
				msgKey,
				zap.String(msgInfo, fmt.Sprintf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)),
				trace,
			)
		}
	}
}

func getTrace(ctx context.Context) string {
	if span := otelTrace.SpanContextFromContext(ctx); span.HasTraceID() {
		return span.TraceID().String()
	}
	return ""
}
