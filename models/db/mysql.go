package db

import (
	"log"
	"m/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DEFAULTDB *gorm.DB

func Initmysql() {
	if db, err := gorm.Open("mysql", config.Uname+":"+config.Password+"@("+config.Host+")/"+config.Db+"?charset=utf8&parseTime=True&loc=Local"); err != nil {
		log.Printf("DEFAULTDB数据库启动异常%S", err)
	} else {
		DEFAULTDB = db
		DEFAULTDB.LogMode(true)
		DEFAULTDB.DB().SetMaxIdleConns(10)
		DEFAULTDB.DB().SetMaxOpenConns(100)
		DEFAULTDB.SingularTable(true)
	}
}
