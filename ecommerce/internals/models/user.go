package models

import "gorm.io/gorm"


type UserModel struct {
	ID           uint      `gorm:"primary key:autoIncrement" json:id`
    Email       *string    `json:email`
	Username    *string    `json:username`  
	Password    *string    `json:password`
}

func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&UserModel{})
    return err	
}

