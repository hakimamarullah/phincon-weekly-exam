package repository

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *sql.DB

func ConnectDatabaseMySQL() {
	dsn := "root:root@tcp(127.0.0.1:3307)/phincon?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	DB, err = db.DB()
	if err != nil {
		panic(err.Error())
	}
}
