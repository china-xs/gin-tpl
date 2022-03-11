// Package db
// @author: xs
// @date: 2022/3/10
// @Description: db,描述,gorm 常规使用、gorm.DB mock DB,当前版本不实现主从配置 会预留配置结构
package db

import (
	"fmt"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type DbOptions struct {
	//  url: products:123456@tcp(db:3306)/products?charset=utf8&parseTime=True&loc=Local
	Host            string        `yaml:"host"`     // ip
	Port            int32         `yaml:"port"`     // 端口
	Database        string        `yaml:"database"` // 表名称
	User            string        `yaml:"user"`     // 用户名
	Pwd             string        `yaml:"pwd"`      // 密码
	PreFix          string        `yaml:"prefix"`   // 表前缀
	MaxIdleConn     int           `yaml:"maxIdleConn"`
	MaxOpenConn     int           `yaml:"maxOpenConn"` // 最大链接池数
	ConnMaxLifetime time.Duration // 单个链接有效时长
	// sources 主库 replicas 从库 稍微大一点的 一主多从即可，当前仅提供多从
	// slaves['从库名称，自定义名称']从库配置
	// db.Clauses(dbresolver.Use("secondary")).First(&user) 指定secondary 库查询
	Slaves                 map[string]string
	Level                  logger.LogLevel `yaml:"level"`
	SlowThreshold          time.Duration   `yaml:"slowTime"`               // 慢查询阀值
	SkipDefaultTransaction bool            `yaml:"skipDefaultTransaction"` // true 开启禁用事物，大约 30%+ 性能提升
}

var ProviderSet = wire.NewSet(New, NewDb)

func New(v *viper.Viper) (*DbOptions, error) {
	var err error
	o := new(DbOptions)
	// 默认全打日志
	o.Level = logger.Info
	// 2秒算慢查询
	o.SlowThreshold = 2 * time.Second
	o.MaxIdleConn = 10
	o.MaxOpenConn = 5
	// 默认配置 8小时，常规使用 1小时
	o.ConnMaxLifetime = 3600
	o.SkipDefaultTransaction = true
	if err = v.UnmarshalKey("db", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}
	return o, err
}

func NewDb(c *DbOptions, log *zap.Logger) (*gorm.DB, func(), error) {
	lOpts := GLOptions{
		Level:         c.Level,
		SlowThreshold: c.SlowThreshold,
	}
	lg := NewGL(&lOpts, log)
	config := gorm.Config{
		SkipDefaultTransaction: c.SkipDefaultTransaction,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.PreFix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,     // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		// 数据库配置 统一公用重写日志库即可
		Logger: lg,
	}
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.User,
		c.Pwd,
		c.Host,
		c.Port,
		c.Database,
	)
	gdb, err := gorm.Open(mysql.Open(dsn), &config)
	if err != nil {
		panic(fmt.Sprintf("init db err:%v", err))
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		panic(fmt.Sprintf("db get mysql.DB err%v", err))
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime * time.Second)
	cleanup := func() {
		sqlDB.Close()
	}
	return gdb, cleanup, nil
}
