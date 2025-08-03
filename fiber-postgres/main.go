package main

import (
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string      `json:"status"`
}

type Repository struct {
	DB *gorm.DB
}

func(r *Repository) CreateTodo(ctx *fiber.Ctx) error {
   todo := Todo {}

   // change the json body into what way the golang undertstands
   err := ctx.BodyParser(&todo)

   if err != nil {
	ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map {"message" : "Request Failed"})
	return err
   }
   
   err = r.DB.Create(&todo).Error

   if err != nil {
	ctx.Status(http.StatusBadRequest).JSON(&fiber.Map {"message" : "could not able to create todo"})
		return err
   }

   ctx.Status(http.StatusOK).JSON(&fiber.Map{ "message" : "Book has been added"})
   return nil
}

func(r *Repository)GetTodos(ctx *fiber.Ctx) error {
	todoModels := &[]models.Todos{}

    err := r.DB.Find(todoModels).Error
	if err != nil {
        ctx.Status(http.StatusBadGateway).JSON(
			&fiber.Map {"message" : "Could not get todos"})
		}
	
    ctx.Status(http.StatusOK).JSON(
		&fiber.Map {"message" : "todos fetched successfully",
	     "books" : todoModels})
	return nil
} 

func(r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/todo",r.CreateTodo)
    api.Delete("/todo/:id",r.DeleteTodo)
	api.Put("/todo/:id",r.UpdateTodo)
	api.Get("/todo/todos",r.GetTodos)
	api.Get("/todo/:id",r.GetTodoById)
}
func main() {
	err := godotenv.Load("env")
    
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}

	r := Repository {
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.listen(":8080")

}