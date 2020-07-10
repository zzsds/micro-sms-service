package service

import (
	"fmt"
	"sync"

	"github.com/zzsds/micro-sms-service/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zzsds/micro-sms-service/conf"
)

var (
	once sync.Once
	Db   *gorm.DB
)

type mysql struct {
}

func InitConnectionMysql() {
	d := NewMysql()
	once.Do(func() {
		if err := d.connection(); err != nil {
			panic(fmt.Sprintf("mysql initmysql failed %s\n", err.Error()))
		}
	})
	d.init()
}

func NewMysql() *mysql {
	return &mysql{}
}

func (m *mysql) connection() (err error) {
	Db, err = gorm.Open("mysql", m.GetDbConnMsg())
	if err != nil {
		return err
	}
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
	Db.SingularTable(true)
	Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8")
	Db.LogMode(conf.Conf.Db.Debug)
	return
}

func (m *mysql) GetDbConnMsg() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", conf.Conf.Db.User, conf.Conf.Db.Password, conf.Conf.Db.Host, conf.Conf.Db.Name, conf.Conf.Db.Charset)
}

func (m *mysql) init() {
	// 数据迁移
	Db.AutoMigrate(
		&models.Send{},
		&models.Code{},
		&models.Template{},
	)
}

func Scope(key string, mode interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(key+" = ?", mode)
	}
}
