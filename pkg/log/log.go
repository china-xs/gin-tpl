// Package log
// @author: xs
// @date: 2022/3/9
// @Description: log,描述
package log

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/wire"
	"github.com/spf13/viper"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const msgTrace = "trace_id"
const msgSpan = "span_id"

const format = "2006-01-02 15:04:05"
const formatFolder = "2006-01-02"

// Options is log configuration struct
type Options struct {
	Filename   string `yaml:"filename"`   // 文件名称
	MaxSize    int    `yaml:"maxSize"`    // 最大文件
	MaxBackups int    `yaml:"maxBackups"` // 最大备份数
	MaxAge     int    `yaml:"maxAge"`     //保留时长天 days
	Level      string `yaml:"level"`      // 日志登记 对应zap.level
	Stdout     bool   `yaml:"stdout"`     // 是否在终端输出
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
func New(o *Options) (*zap.Logger, func(), error) {
	var (
		err    error
		level  = zap.NewAtomicLevel()
		logger *zap.Logger
	)

	err = level.UnmarshalText([]byte(o.Level))
	if err != nil {
		return nil, func() {}, err
	}
	write := &lumberjack.Logger{ // concurrent-safed
		Filename:   o.Filename,   // 文件路径
		MaxSize:    o.MaxSize,    // MaxSize 兆字节
		MaxBackups: o.MaxBackups, // 最多保留 300 个备份
		MaxAge:     o.MaxAge,     // 最大时间，默认单位 day
		LocalTime:  true,         // 使用本地时间
		Compress:   false,        // 是否压缩 disabled by default
	}

	fw := zapcore.AddSync(write)
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
	fc := func() {
		logger.Sync() // 缓存
		write.Close() // os close
	}
	return logger, fc, err
}

type Log struct {
	l *zap.Logger
}

func NewL(l *zap.Logger) *Log {
	return &Log{
		l: l,
	}
}

func (l Log) With(ctx context.Context) *zap.Logger {
	var traceId, spanId string
	if span := otelTrace.SpanContextFromContext(ctx); span.HasTraceID() {
		traceId = span.TraceID().String()
	}
	if span := otelTrace.SpanContextFromContext(ctx); span.HasSpanID() {
		spanId = span.SpanID().String()
	}
	// 跳过 gorm 内置的调用
	var (
		gormPackage = filepath.Join("gorm.io", "gorm")
		genPackage  = filepath.Join("gorm.io", "gen")

		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)
	var fields []zap.Field
	fields = append(fields,
		zap.String(msgTrace, traceId),
		zap.String(msgSpan, spanId),
		zap.String("caller", GetCaller()),
	)
	// 减去一次封装，以及一次在 logger 初始化里添加 zap.AddCallerSkip(1)
	clone := l.l.WithOptions(zap.AddCallerSkip(-2))
	for i := 2; i < 20; i++ {
		_, file, _, ok := runtime.Caller(i)
		fmt.Println("file:", file)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, genPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			// 返回一个附带跳过行号的新的 zap logger
			return clone.WithOptions(zap.AddCallerSkip(i), zap.Fields(fields...))
		}
	}
	fmt.Println("?")
	return l.l.With(fields...)
}

func WithCtx(ctx context.Context, log *zap.Logger) *zap.Logger {
	var traceId, spanId string
	if span := otelTrace.SpanContextFromContext(ctx); span.HasTraceID() {
		traceId = span.TraceID().String()
	}
	if span := otelTrace.SpanContextFromContext(ctx); span.HasSpanID() {
		spanId = span.SpanID().String()
	}
	var fields []zap.Field
	fields = append(fields,
		zap.String(msgTrace, traceId),
		zap.String(msgSpan, spanId),
		zap.String("caller", GetCaller()),
	)
	return log.With(fields...)
}

func fileWire() (io.Writer, func()) {
	l := &lumberjack.Logger{ // concurrent-safed
		Filename:   "app.log",   // 文件路径
		MaxSize:    1024 * 1024, // 1T // MaxSize 不设置单个文件最大为100M
		MaxBackups: 0,           // 最多保留 300 个备份
		MaxAge:     365,         // 最大时间，默认单位 day
		LocalTime:  true,        // 使用本地时间
		Compress:   false,       // 是否压缩 disabled by default
	}
	return l, func() { l.Close() }
}

func GetCaller() string {
	depth := 3
	_, file, line, _ := runtime.Caller(depth)
	//fmt.Printf("file:%v\n",file)
	// gorm db 回调层
	if strings.HasSuffix(file, "callbacks.go") {
		return ""
		//depth++
		//_, file, line, _ = runtime.Caller(7)
	}
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
