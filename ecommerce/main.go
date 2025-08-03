package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {

	app := fiber.New()
    app.Listen(":8080")
}