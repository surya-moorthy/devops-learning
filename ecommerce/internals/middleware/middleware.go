package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)


var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func JWTMiddleware(ctx *fiber.Ctx) error {
	
	authHeader := ctx.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
	}

    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
				}
			}

		ctx.Locals("user", claims)
		
	return ctx.Next()
} 