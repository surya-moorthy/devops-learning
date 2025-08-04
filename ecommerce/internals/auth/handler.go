package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/devops-learning/ecommerce/internals/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	service *AuthRepository
}

var secretkey = []byte(os.Getenv("SECRET_KEY"))
func (h *AuthRepository) Login(context *fiber.Ctx) error {
	// Define input struct
	user := models.User{}

	// Parse the request body
	if err := context.BodyParser(&user); err != nil {
		return context.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	// Get user from DB using email
	getUser, err := h.GetUserByEmail(&user.Email)
	if err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// log.Fatal("user password:",user)
	// log.Fatal("db user password:",getUser.Password)
	// log.Fatal("user password:",user)

	// Compare stored hash with incoming password
	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(user.Password)); err != nil {
		return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	// Generate JWT token
	token, err := createToken(getUser.Username, getUser.Email)
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create a token",
		})
	}

	// Respond with success and token
	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged in successfully",
		"token":   token,
	})
}


func(h *AuthRepository) Register(context *fiber.Ctx) error {
    user := models.User {}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map {"message" : "Request Failed"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return context.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to process new password",
		})
	}
	
	// Update user with new password and clear reset token
	user.Password = string(hashedPassword)

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

func (h *AuthRepository) ResetPassword(ctx *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	
	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}
	
	// Validate required fields
	if req.Email == "" || req.NewPassword == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Email, token, and new password are required",
		})
	}
	
	// Validate password strength
	if len(req.NewPassword) < 8 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 8 characters long",
		})
	}
	
	// Find user by email
	foundUser, err := h.FindUserByEmail(&req.Email)
	if err != nil || foundUser == nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reset request",
		})
	}
	
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to process new password",
		})
	}
	
	// Update user with new password and clear reset token
	foundUser.Password = string(hashedPassword)
	
	if err := h.UpdateUser(foundUser); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update password",
		})
	}
	
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

// // Helper function to generate secure token
// func generateSecureToken() (string, error) {
// 	bytes := make([]byte, 32)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(bytes), nil
// }
// func(h *AuthRepository) Logout(context *fiber.Ctx) error {
   
// 	context.Status(http.StatusOK).JSON( &fiber.Map{"message" : "User Logged out Successfully"})
// 	return nil
// }

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