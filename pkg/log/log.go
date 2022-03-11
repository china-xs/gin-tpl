// Package log
// @author: xs
// @date: 2022/3/9
// @Description: log,描述
package log

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/google/wire"
	"github.com/spf13/viper"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const msgTrace = "trace_id"

// Options is log configuration struct
type Options struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      string
	Stdout     bool
}

func NewOptions(v *viper.Viper) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("log", o); err != nil {
		return nil, err
	}

	return o, err
}

// New for init zap log library
func New(o *Options) (*zap.Logger, error) {
	var (
		err    error
		level  = zap.NewAtomicLevel()
		logger *zap.Logger
	)

	err = level.UnmarshalText([]byte(o.Level))
	if err != nil {
		return nil, err
	}

	fw := zapcore.AddSync(&lumberjack.Logger{
		Filename:   o.Filename,
		MaxSize:    o.MaxSize, // megabytes
		MaxBackups: o.MaxBackups,
		MaxAge:     o.MaxAge, // days
	})

	cw := zapcore.Lock(os.Stdout)

	// file core 采用jsonEncoder
	cores := make([]zapcore.Core, 0, 2)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		//CallerKey:     "caller",// kratos 已经配置 caller zap 负责写入数据即可
		//EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}
	je := zapcore.NewJSONEncoder(encoderConfig)
	cores = append(cores, zapcore.NewCore(je, fw, level))

	// stdout core 采用 ConsoleEncoder
	if o.Stdout {
		ce := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		cores = append(cores, zapcore.NewCore(ce, cw, level))
	}

	core := zapcore.NewTee(cores...)
	logger = zap.New(core)

	zap.ReplaceGlobals(logger)

	return logger, err
}

//
// WithCtx
// @Description: 根据上下文
// @param ctx
// @return zap.Field
//
func WithCtx(ctx context.Context) []zap.Field {
	var trace string
	if span := otelTrace.SpanContextFromContext(ctx); span.HasTraceID() {
		trace = span.TraceID().String()
	}
	var fields []zap.Field
	fields = append(fields,
		zap.String(msgTrace, trace),
		zap.String("caller", GetCaller()),
	)

	return fields
}

func GetCaller() string {
	depth := 3
	_, file, line, _ := runtime.Caller(depth)
	//fmt.Printf("file:%v\n",file)
	// 处理数据库层
	//if strings.HasSuffix(file, "dbLogger.go") {
	//	return ""
	//	//depth++
	//	//_, file, line, _ = runtime.Caller(7)
	//}
	idx := strings.LastIndexByte(file, '/')
	//fmt.Printf("caller:%v\n",idx)
	return file[idx+1:] + ":" + strconv.Itoa(line)
}

func RestyLog(resp *resty.Response, field ...zap.Field) []zap.Field {
	traceInfo := resp.Request.TraceInfo()
	field = append(field,
		zap.String("url", resp.Request.URL),
		zap.Int("resp_status_code", resp.StatusCode()),
		zap.String("resp_status", resp.Status()),
		zap.String("resp_time", resp.Time().String()),
		zap.String("resp_received", resp.ReceivedAt().GoString()),
		zap.String("request_DNSLookup", traceInfo.DNSLookup.String()),
		zap.String("request_ConnTime", traceInfo.ConnTime.String()),
		zap.String("request_TCPConnTime", traceInfo.TCPConnTime.String()),
		zap.String("request_TLSHandshake", traceInfo.TLSHandshake.String()),
		zap.String("request_ServerTime", traceInfo.ServerTime.String()),
		zap.String("request_ResponseTime", traceInfo.ResponseTime.String()),
		zap.String("request_TotalTime", traceInfo.TotalTime.String()),
		zap.Bool("request_IsConnReused", traceInfo.IsConnReused),
		zap.Bool("request_IsConnWasIdle", traceInfo.IsConnWasIdle),
		zap.String("request_ConnIdleTime", traceInfo.ConnIdleTime.String()),
	)
	return field
}

var ProviderSet = wire.NewSet(New, NewOptions)
