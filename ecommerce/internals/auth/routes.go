package auth

import (
	"github.com/devops-learning/ecommerce/internals/middleware"
	"github.com/gofiber/fiber/v2"
)

func(r *AuthRepository)  SetupAuthRoutes(app fiber.Router) {
	app.Post("/register",r.Register)
	app.Post("/login",r.Login)
	auth := app.Group("/auth",middleware.JWTMiddleware)	
	auth.Post("/logout",r.Logout)
	auth.Post("/reset-password",r.ResetPassword)
}