// Package gorm
// @author: xs
// @date: 2022/3/10
// @Description: gorm
package main

import (
	"fmt"
	"github.com/china-xs/gin-tpl/pkg/config"
	db2 "github.com/china-xs/gin-tpl/pkg/db"
	"go.uber.org/zap"
	"gorm.io/gen"
)

func main() {

	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../internal/data/dao/query",
		FieldNullable: true,
	})
	l, _ := zap.NewProduction()
	// config 注意表前缀问题
	v, err := config.New("../../configs/app.yaml")
	opts, err := db2.New(v)
	db, fc, err := db2.NewDb(opts, l)
	if err != nil {
		panic(fmt.Sprintf("init-db:%v", err.Error()))
	}
	defer fc()
	// 复用工程原本使用的SQL连接配置
	g.UseDB(db)
	// 所有需要实现查询方法的结构体 增加表不能把原来表删除
	g.ApplyBasic(
		//g.GenerateModel("t_mp_list", gen.FieldType("state", "int")),
		g.GenerateModel("oa_departments"),
	)
	g.Execute()

}
