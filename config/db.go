package config

import (
    "aurora-im/model"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func InitDB() {
    dsn := "aurora_user:aurora_pass@tcp(127.0.0.1:3306)/aurora_db?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    // Assign global DB variable
    DB = db

    // Auto migrate tables
    err = DB.AutoMigrate(&model.User{})
    if err != nil {
        log.Fatal("Failed to migrate user table:", err)
    }
}
