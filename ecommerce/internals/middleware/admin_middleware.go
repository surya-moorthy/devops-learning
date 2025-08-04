package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(ctx *fiber.Ctx) error {
	userClaims, ok := ctx.Locals("user").(map[string]interface{})
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	role, ok := userClaims["role"].(string)
	if !ok || role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
	}

	return ctx.Next()
}
