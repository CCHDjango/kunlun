package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type dbQuote struct {
	DBHuobi *gorm.DB
}

var DB dbQuote

func init() {
	DB = dbQuote{}
	var mysql = fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"root", "kunlun2020", "127.0.0.1", "kunlun_quote",
	)

	db, err := gorm.Open("mysql", mysql)
	if err != nil {
		panic(err.Error())
	}
	//连接池设置
	db.DB().SetMaxIdleConns(6)
	db.DB().SetMaxOpenConns(50)
	DB.DBHuobi = db
}
