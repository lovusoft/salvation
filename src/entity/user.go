package entity

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Name     string `gorm:"UNIQUE"`
	Password string
	Coin     int
	Secrets  []Secret
}

func (u User) UserDel(user User) {
	db, err := gorm.Open("mysql", "lovu:1314@/salvation?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库出错")
	}
	defer db.Close()
	db.Delete(&user)
}

func (u User) UserFind(user User) User {
	db, err := gorm.Open("mysql", "lovu:1314@/salvation?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库出错")
	}
	defer db.Close()
	db.First(&user)
	return user
}
