package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open("mysql","root:@(localhost)/db_explore_gin?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("Connection Failed")
	} 

	db.AutoMigrate(&Student{})

	return db
}