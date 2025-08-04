package auth

import "github.com/gofiber/fiber/v2"

func(r *AuthRepository)  SetupAuthRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Post("/register",r.Register)
	auth.Post("/login",r.Login)
	auth.Post("/logout",r.Logout)
	auth.Post("/reset-password",r.ResetPassword)
}