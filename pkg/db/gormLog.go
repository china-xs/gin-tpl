// Package db
// @author: xs
// @date: 2022/3/9
// @Description: log 重写gorm logger 日志模块 仅写操作日志
package db

import (
	"context"
	"fmt"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

const msgKey = "sql"
const msgInfo = "sql-info"
const msgTrace = "trace_id"

var ProviderGLSet = wire.NewSet(NewGL, NewGLOpts)

var (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

type GLog struct {
	zapLog        *zap.Logger
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration // 慢查询阀值
}

type GLOptions struct {
	Level         logger.LogLevel `yaml:"level"`
	SlowThreshold time.Duration   `yaml:"slowTime"` // 慢查询阀值
}

func NewGLOpts(v *viper.Viper) (*GLOptions, error) {
	var (
		err error
		o   = new(GLOptions)
	)
	if err = v.UnmarshalKey("db", o); err != nil {
		return nil, err
	}

	return o, err
}
func NewGL(o *GLOptions, l *zap.Logger) *GLog {
	return &GLog{
		zapLog:        l,
		LogLevel:      o.Level,
		SlowThreshold: o.SlowThreshold,
	}
}

// NewDBLog 依赖
func NewDBLog(l *zap.Logger) *GLog {
	return &GLog{
		zapLog:        l,
		LogLevel:      logger.Info,     //打印所有日志，登记越小打印越少
		SlowThreshold: 2 * time.Second, // 2秒
	}
}

// LogMode log mode
func (l *GLog) LogMode(level logger.LogLevel) logger.Interface {
	newLog := *l
	newLog.LogLevel = level
	return &newLog
}

// Info print info
func (l GLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		var fields []zap.Field
		fields = []zap.Field{
			zap.String(msgInfo, fmt.Sprintf(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
		}
		fields = append(fields, log.WithCtx(ctx)...)
		l.zapLog.Info(msgKey, fields...)
	}
}

// Warn print warn messages
func (l GLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		var fields []zap.Field
		fields = []zap.Field{
			zap.String(msgInfo, fmt.Sprintf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
		}
		fields = append(fields, log.WithCtx(ctx)...)
		l.zapLog.Warn(msgKey, fields...)
	}
}

// Error print error messages
func (l GLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		var fields []zap.Field
		fields = []zap.Field{
			zap.String(msgInfo, fmt.Sprintf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)),
		}
		fields = append(fields, log.WithCtx(ctx)...)
		l.zapLog.Error(msgKey, fields...)
	}
}

// Trace print sql message
func (l GLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	var fields []zap.Field
	fields = append(fields, log.WithCtx(ctx)...)
	elapsed := time.Since(begin)
	switch {
	//|| !l.IgnoreRecordNotFoundError
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			fields = append(fields, zap.String(msgInfo, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)))
		} else {
			fields = append(fields, zap.String(msgInfo, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)))
		}
		l.zapLog.Error(msgKey, fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			fields = append(fields, zap.String(msgInfo, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)))
		} else {
			fields = append(fields, zap.String(msgInfo, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)))
		}
		l.zapLog.Warn(msgKey, fields...)
	case l.LogLevel == logger.Info:
		sql, rows := fc()

		if rows == -1 {
			fields = append(fields, zap.String(msgInfo, fmt.Sprintf(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)))
		} else {
			fields = append(fields, zap.String(msgInfo, fmt.Sprintf(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)))
		}
		l.zapLog.Info(msgKey, fields...)
	}
}
