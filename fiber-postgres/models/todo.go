package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	ID           uint        `gorm:"primary key;autoIncrement" json:"id"`
	Title        *string      `json:"title"`
	Description  *string      `json:description`
	Status       *string      `json:status`
}

func MigrateTodo(db *gorm.DB) error {
   err := db.AutoMigrate(&Todo{})
   return err

}