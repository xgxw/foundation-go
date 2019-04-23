package mysql

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

// Options 创建数据库的选项
type Options struct {
	Driver    string `yaml:"driver" mapstructure:"driver"`
	Dsn       string `yaml:"dsn" mapstructure:"dsn"`
	KeepAlive int    `yaml:"keep_alive" mapstructure:"keep_alive"`
	MaxIdles  int    `yaml:"max_idles" mapstructure:"max_idles"`
	MaxOpens  int    `yaml:"max_opens" mapstructure:"max_opens"`
}

// DB Gorm封装
type DB struct {
	*gorm.DB

	// ticker 用于keep alive的定时器
	ticker *time.Ticker
	// SavePoint
	txDepth int
}

// ErrDBRecordNotFound 未查询到数据库记录
var ErrDBRecordNotFound = gorm.ErrRecordNotFound

// NewDatabase 创建新的数据库对象
func NewDatabase(opts Options) (*DB, error) {
	o, err := gorm.Open(opts.Driver, opts.Dsn)
	if err != nil {
		return nil, errors.Wrap(err, "database open failed")
	}

	db := &DB{DB: o}

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

func (db *DB) keepAlive(d time.Duration) {
	db.ticker = time.NewTicker(d)
	go func() {
		for range db.ticker.C {
			db.DB.DB().Ping()
		}
	}()
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	if db.ticker != nil {
		db.ticker.Stop()
	}
	return db.DB.Close()
}
