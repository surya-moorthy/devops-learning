package main

import (
	"log"
	"os"

	"github.com/devops-learning/ecommerce/internals/auth"
	"github.com/devops-learning/ecommerce/internals/db"
	"github.com/devops-learning/ecommerce/internals/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {

	 config := &db.Config {
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		DBname: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	dbconn , err := db.NewConnection(config) 
	
	if err != nil {
		log.Fatal("could not load the database")
	}

	err = models.MigrateAll(dbconn)

	if err != nil {
		log.Fatal("could not migrate db")
	}

	authrepo := auth.AuthRepository{
		DB: dbconn,
	}

	app := fiber.New()
    app.Listen(":8080")
}