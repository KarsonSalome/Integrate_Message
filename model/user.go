package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
    Username string `gorm:"unique;not null" json:"username"`
	Phone string `gorm:"unique;not null" json:"phone"`
    Password string `gorm:"not null" json:"password"`
	Avatar string `gorm:"not null" json:"avatar"`
}
