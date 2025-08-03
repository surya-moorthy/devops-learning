package main

import "github.com/gin-gonic/gin"


type PostTask struct {
	Title       string `json:"task" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=progress success"`
}

type Task struct {
	Title       string `json:"task"`
	Description  string `json:"description"`
	Status      string `json:"status"`
}

var task Task

func main() {
	server := gin.Default();  // helps for logging funcs andrecoveries

	server.GET("/",func(ctx *gin.Context) {
		ctx.JSON(200,gin.H {
			"task" : task,
		})
	})
	server.POST("/",postTask)
	server.PUT("/:id",updateTask)


	server.Run(":8080")
}

func postTask(ctx *gin.Context) {
   var input PostTask;

   if err := ctx.ShouldBindJSON(&input); err != nil {
	  ctx.JSON(400 , gin.H {"error" : err.Error(),})
	  return 
   }

  if (input.Title == "" || input.Description == "") {
	ctx.JSON(404,gin.H {"message" : "Provide valid input",})
	return
  } 

  task = Task{
	Title:       input.Title,
	Description: input.Description,
	Status:      input.Status,
  }

  ctx.JSON(200,gin.H {
	"message" : "data has been posted",
	"task" : input,
  })
}

func updateTask(ctx *gin.Context) {
	var input Task;

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400 , gin.H {"error" : err.Error(),})
	   return 
	}

	// Apply only the fields that are provided
	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Status != "" {
		task.Status = input.Status
	}


	ctx.JSON(200,gin.H {
	"message" : "data has been updated",
	"task" : task,
  })
}

