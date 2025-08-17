package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/devops-learning/fiber-postgres/models"
	"github.com/devops-learning/fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Todo struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
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
	todoModels := &[]models.Todo{}

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

func(r *Repository) DeleteTodo(ctx *fiber.Ctx) error {
   todoModels := models.Todo{}
   id := ctx.Params("id")
   if id == "" {
	ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
		"message":"id cannot be empty",
	})
	return nil
   }

   err := r.DB.Delete(todoModels,id); 
   if err.Error != nil {
	ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message" : "could not delete todo",
	})
	return err.Error
   }

   ctx.Status(http.StatusOK).JSON(&fiber.Map{
	"message": "todo deleted successfully",
   })
   return nil
}

func(r *Repository) UpdateTodo(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    todo := Todo {}

   // change the json body into what way the golang undertstands
   err := ctx.BodyParser(&todo)

   if err != nil {
	ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map {"message" : "Request Failed"})
	return err
   }

	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map {"message" : "id cannot be empty",})
		return nil
	}

	err = r.DB.Model(&models.Todo{}).Where("id = ?", id).Updates(&todo).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not update todo",
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "todo updated successfully",
	})
	return nil
}

func(r *Repository) GetTodoById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	bookModel := models.Todo{}

	fmt.Println("the id req ",id)

	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map {"message" : "id cannot be empty",})
		return nil
	}

	fmt.Println("the id is ", id)

	err := r.DB.Where("id = ?", id).First(&bookModel).Error

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message" : "could not able to delete the todo",
		})
		return  nil
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map {
		"message": "todo with the given id fetched successfully",
		"data" : bookModel,
	})
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
	err := godotenv.Load(".env")
    
	if err != nil {
		log.Fatal(err)
	}

    config := &storage.Config {
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		DBname: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}


	err = models.MigrateTodo(db)

	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository {
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}