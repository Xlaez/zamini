package middlewares

import (
	"ecommerce/pkg/helpers"

	"github.com/gofiber/fiber/v2"
)

func RequireAuthMiddleware(c *fiber.Ctx) error{
	authHeader := c.Get("Authorization")
	token := c.Cookies("x-auth-jwt")

	if authHeader == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"message": "Unauthroized, try singing in",
		})
	}

	if token == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Try to signin first",
		})
	}

	id, userType, err := helpers.VerifyToken(token)

	if err !=nil{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
		})
	}

	c.Locals("id", id)
	c.Locals("userType", userType)
	return c.Next()
}