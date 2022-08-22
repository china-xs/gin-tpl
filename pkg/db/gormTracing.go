// Package db
// @author: xs
// @date: 2022/6/29
// @Description: db
package db

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"

	gormSpanKey = "__gorm_span"
)

// 告诉编译器这个结构体实现了gorm.Plugin接口
var _ gorm.Plugin = &GormTrace{}

type GormTrace struct {
}

func (op *GormTrace) Name() string {
	return "opentracingPlugin"
}

func (op *GormTrace) Initialize(db *gorm.DB) error {
	//执行前
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 执行结束
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return nil
}

func before(db *gorm.DB) {
	ctx := db.Statement.Context
	ctx, span := otel.GetTracerProvider().
		Tracer("gorm.io").
		Start(ctx, "query-start")
	db.InstanceSet(gormSpanKey, span)
	db.Statement.Context = ctx
}
func after(db *gorm.DB) {
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		// 不存在就直接抛弃掉
		return
	}
	span, ok := _span.(trace.Span)
	if !ok {
		return
	}
	span.End()

}
