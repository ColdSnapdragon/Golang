package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// dao(Data Access Object,数据操作对象)负责与数据库进行交互，对数据进行持久化操作

var (
	DB *gorm.DB
)

func ConnectMySQL() (err error) {
	dsn := "root:380024507@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}
