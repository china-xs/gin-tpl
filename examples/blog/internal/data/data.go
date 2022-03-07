// Package data
// @author: xs
// @date: 2022/3/7
// @Description: data db、cache缓存还没考虑好直接开多一层还是跟db保持同级
// @Note 返回DB 主要为了写测试用例，使用mock 替代真是db 思考 left join db
package data

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

//
type dbConfig struct {
	//  url: products:123456@tcp(db:3306)/products?charset=utf8&parseTime=True&loc=Local
	Host            string        `yaml:"host"`     // ip
	Port            string        `yaml:"port"`     // 端口
	Database        string        `yaml:"database"` // 表名称
	User            string        `yaml:"user"`     // 用户名
	Pwd             string        `yaml:"pwd"`      // 密码
	PreFix          string        `yaml:"prefix"`
	MaxIdleConn     int           `yaml:"maxIdleConn"`
	MaxOpenConn     int           `yaml:"maxOpenConn"` // 最大链接池数
	ConnMaxLifetime time.Duration // 单个链接有效时长
}
type db struct {
	gdb *gorm.DB
}
type ORM interface {
	DB() *gorm.DB
}

var _ ORM = (*db)(nil)

//logger *zap.Logger
func New(v *viper.Viper) (*dbConfig, error) {
	var err error
	o := new(dbConfig)
	if err = v.UnmarshalKey("db", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}

	//logger.Info("load database options success", zap.String("url", o.URL))

	return o, err
}

func NewDb(c *dbConfig, log *zap.Logger) (ORM, func(), error) {
	dbconfig := gorm.Config{
		//SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		// 数据库配置 统一公用重写日志库即可
		//Logger: NewGLogger(),
	}
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.User,
		c.Pwd,
		c.Host,
		c.Port,
		c.Database,
	)
	gdb, err := gorm.Open(mysql.Open(dsn), &dbconfig)
	if err != nil {
		panic(fmt.Sprintf("init db err:%v", err))
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		panic(fmt.Sprintf("db get mysql.DB err%v", err))
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime * time.Minute)
	cleanup := func() {
		sqlDB.Close()
	}
	return &db{
		gdb: gdb,
	}, cleanup, nil
}

func (this db) DB() *gorm.DB {
	return this.gdb
}
