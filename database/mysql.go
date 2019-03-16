package database

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

// MysqlOptions 创建数据库的选项
type MysqlOptions struct {
	Driver    string `yaml:"driver" mapstructure:"driver"`
	Dsn       string `yaml:"dsn" mapstructure:"dsn"`
	KeepAlive int    `yaml:"keep_alive" mapstructure:"keep_alive"`
	MaxIdles  int    `yaml:"max_idles" mapstructure:"max_idles"`
	MaxOpens  int    `yaml:"max_opens" mapstructure:"max_opens"`
}

// MysqlDB Gorm封装
type MysqlDB struct {
	*gorm.DB

	// ticker 用于keep alive的定时器
	ticker *time.Ticker
}

// ErrDBRecordNotFound 未查询到数据库记录
var ErrDBRecordNotFound = gorm.ErrRecordNotFound

// NewMysqlDatabase 创建新的数据库对象
func NewMysqlDatabase(opts MysqlOptions) (*MysqlDB, error) {
	o, err := gorm.Open(opts.Driver, opts.Dsn)
	if err != nil {
		return nil, errors.Wrap(err, "database open failed")
	}

	db := &MysqlDB{DB: o}

	if opts.MaxIdles > 0 {
		o.DB().SetMaxIdleConns(opts.MaxIdles)
	}
	if opts.MaxOpens > 0 {
		o.DB().SetMaxOpenConns(opts.MaxOpens)
	}
	if opts.KeepAlive > 0 {
		db.keepAlive(time.Second * time.Duration(opts.KeepAlive))
	}

	return db, nil
}

func (db *MysqlDB) keepAlive(d time.Duration) {
	db.ticker = time.NewTicker(d)
	go func() {
		for range db.ticker.C {
			db.DB.DB().Ping()
		}
	}()
}

// Close 关闭数据库连接
func (db *MysqlDB) Close() error {
	if db.ticker != nil {
		db.ticker.Stop()
	}
	return db.DB.Close()
}
