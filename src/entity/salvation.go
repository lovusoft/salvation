package entity

import "github.com/jinzhu/gorm"

type Salvation struct {
	gorm.Model
	Writer      uint
	Adjudicator uint
	Adviser     uint
	Secrets     []Secret
}
