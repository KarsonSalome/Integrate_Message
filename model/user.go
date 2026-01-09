package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
    Username string `gorm:"unique;not null" json:"username"`
	Phone string `gorm:"unique;not null" json:"phone"`
    Password string `gorm:"not null" json:"password"`
	Avatar string `gorm:"not null" json:"avatar"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
    Phone    string `json:"phone" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginRes struct {
	UserID uint `json:"userid"`
    Username string `json:"username"`
    Phone string `json:"phone"`
    Avatar string `json:"avatar"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
    Password string `json:"password"`
}
