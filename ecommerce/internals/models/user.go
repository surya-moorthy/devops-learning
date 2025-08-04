package models

import "gorm.io/gorm"


type UserModel struct {
	ID           uint      `gorm:"primary key:autoIncrement" json:id`
    Email       *string    `gorm:"not null" json:email`
	Username    *string    `json:username`  
	Password    *string    `gorm:"not null" json:password`
}

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(&UserModel{})
    return err	
}

