package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSLMode  string
}

func NewConnection(config *Config)(*gorm.DB,error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s port=%s password=%s dbname=%s sslmode=%s",
		config.Host,config.User,config.Password,config.DBname,config.SSLMode)
	db,err := gorm.Open(postgres.Open(dsn),&gorm.Config{})

	if err != nil {
		return db , err
	}
	return db , nil
}


 