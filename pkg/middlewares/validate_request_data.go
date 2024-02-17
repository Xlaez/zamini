package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateMiddleware(c *fiber.Ctx) error{
	validate := validator.New()

	type UserInputCred struct{
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var user UserInputCred
	if err := c.BodyParser(&user); err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "missing field",
		})
	}

	if err := validate.Struct(&user); err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid request body",
			"err": err.Error(),
		})
	}

	return c.Next()
}