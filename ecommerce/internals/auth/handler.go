package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/devops-learning/ecommerce/internals/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type handler struct {
	service *AuthRepository
}

var secretkey = []byte(os.Getenv("SECRET_KEY"))

func(h *AuthRepository) Login(context *fiber.Ctx) error {
 
	user := models.UserModel{}

	err := context.BodyParser(&user)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{ "message" : "Request Failed"},
		)
	}

    getUser , err := h.GetUserByEmail(user.Email,user.Password)

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{ "message" : "User cannot found"},
		)
	}

    token , err := createToken(*getUser.Username,*getUser.Email)

	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map {
				"message" : "Failed to create a token",
			})
	}

	context.Status(http.StatusOK).JSON( &fiber.Map{
		"message" : "User Logged in Successfully",
	    "token" : token,
	})
	return nil
}

func(h *AuthRepository) Register(context *fiber.Ctx) error {
    user := models.UserModel {}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map {"message" : "Request Failed"})
	}

	err = h.CreateUser(&user)
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Failed to create user",
		   "error" : err.Error(),
		})
	}

	context.Status(http.StatusOK).JSON( &fiber.Map{"message" : "User created Successfully"})
	return nil
}

func(h *AuthRepository) ResetPassword(context *fiber.Ctx) error {
    user := models.UserModel {}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map {"message" : "Request Failed"})
	}

	err = h.FindUserByEmail(user.Email)

    if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map {"message" : "User not Found"})
	}
    
	err = h.UpdateUser(&user)

    if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map {"message" : "failed to update the password"})
	}
	
   return nil
}

func(h *AuthRepository) Logout(context *fiber.Ctx) error {
   
	context.Status(http.StatusOK).JSON( &fiber.Map{"message" : "User Logged out Successfully"})
	return nil
}


func createToken(username string,email string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
		"email" : email,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })
    tokenString, err := token.SignedString(secretkey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}