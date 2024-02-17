package routes

import (
	"ecommerce/pkg/controllers"
	"ecommerce/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App){
	api := app.Group("/api/v1")

	userApi := api.Group("/auth")

	userApi.Post("/signup", middlewares.ValidateMiddleware ,controllers.Signup)
}