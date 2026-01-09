package model

import "gorm.io/gorm"

type Contact struct {
    gorm.Model

    OwnerID   uint `gorm:"not null"` // The user who added the contact
    ContactID uint `gorm:"not null"` // The person being added

    LastMsg string `gorm:"type:varchar(100)" json:"lastmsg"` 
    // "added" = normal contact, "blocked", "favorite", etc.

    // Relations
    Owner   User `gorm:"foreignKey:OwnerID;references:ID"`   // owner info
    Contact User `gorm:"foreignKey:ContactID;references:ID"` // contact info
}

type SearchRes struct {
	UserID uint `json:"userid"`
    Username string `json:"username"`
    Phone string `json:"phone"`
    Avatar string `json:"avatar"`
}

type SearchReq struct {
    Phone string `json:"phone"`
}

type AddContactReq struct {
    ContactID uint `json:"contact_id"`
}