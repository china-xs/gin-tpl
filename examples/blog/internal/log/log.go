// Package log
// @author: xs
// @date: 2022/3/7
// @Description: log 暂时干不掉整个框架依赖log 特殊处理问题 根据上下文写入日志
package log

import (
	"context"
	"fmt"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const traceKey = "trace_id"
const callerKey = "caller"
const timeFormat = "2006-01-02"

func NewLog() *zap.Logger {
	l, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("初始化zap.logger,err:%v", err.Error()))
	}
	// 直接保存接口
	//ctx := context.TODO()
	//l.With(WithCtx(ctx)...)
	return l
}

func WithCtx(ctx context.Context) []zap.Field {
	var trace string
	if span := otelTrace.SpanContextFromContext(ctx); span.HasTraceID() {
		trace = span.TraceID().String()
	}
	var fields []zap.Field

	fields = append(fields,
		zap.String(traceKey, trace),
		zap.String(callerKey, getCaller()),
	)
	return fields
}

// depth 层数有待检查
func getCaller() string {
	depth := 3
	_, file, line, _ := runtime.Caller(depth)
	//fmt.Printf("file:%v\n",file)
	// 处理数据库层
	if strings.HasSuffix(file, "dbLogger.go") {
		return ""
		//depth++
		//_, file, line, _ = runtime.Caller(7)
	}
	idx := strings.LastIndexByte(file, '/')
	//fmt.Printf("caller:%v\n",idx)
	return file[idx+1:] + ":" + strconv.Itoa(line)
}

//按天分割日志
func withFile(path, filename string) io.Writer {
	file := getFilePath(path, filename)
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}
	l := &lumberjack.Logger{ // concurrent-safed
		Filename:   file,        // 文件路径
		MaxSize:    1024 * 1024, // 1T // MaxSize 不设置单个文件最大为100M
		MaxBackups: 0,           // 最多保留 300 个备份
		MaxAge:     0,           // 最大时间，默认单位 day
		LocalTime:  true,        // 使用本地时间
		Compress:   false,       // 是否压缩 disabled by default
	}
	go func() {
		for {
			<-time.After(time.Hour * 24)
			time.Now().Format("")
			// 重写文件路径
			l.Filename = getFilePath(path, filename)
			l.Rotate()
		}
	}()
	//ticker := time.NewTicker(10 * time.Second)
	//go func(t *time.Ticker) {
	//	for {
	//		<-t.C
	//		if t := time.Now().Format(DefaultTimeLayoutDay); t != dayTime {
	//			l.Filename = getFilePath(path, filename)
	//			l.Rotate()
	//			dayTime = t
	//		}
	//	}
	//}(ticker)
	return l
}

func getFilePath(path, filename string) string {
	dayTime := time.Now().Format("2006-01-02")
	return path + dayTime + "/" + filename
}
