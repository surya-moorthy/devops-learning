package main

import (
	"log"
	"os"

	"github.com/devops-learning/ecommerce/internals/auth"
	"github.com/devops-learning/ecommerce/internals/db"
	"github.com/devops-learning/ecommerce/internals/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
    

	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")


	 config := &db.Config {
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		User: os.Getenv("POSTGRES_USER"),
		DBname: os.Getenv("POSTGRES_DB"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	dbconn , err := db.NewConnection(config) 
	
	if err != nil {
		log.Fatal("could not load the database")
		log.Fatal(err)
	}

	err = models.MigrateAll(dbconn)

	if err != nil {
		log.Fatal("could not migrate db")
	}

	authrepo := auth.AuthRepository{
		DB: dbconn,
	}

	authrepo.SetupAuthRoutes(v1)
	//userrepo.SetupUserRoutes(app,middleware)
	//productrepo.SetupProductRoutes(app,middleware)

    app.Listen(":8080")
}