package routes

import (
	"ecommerce/pkg/controllers"
	"ecommerce/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App){
	api := app.Group("/api/v1")

	authApi := api.Group("/auth")

	userApi := api.Group("/user")

	authApi.Post("/signup", middlewares.ValidateMiddleware ,controllers.Signup)
	authApi.Post("/signin", middlewares.ValidateMiddleware, controllers.SingIn)
	authApi.Post("/signout", middlewares.RequireAuthMiddleware, controllers.SingOut)
	userApi.Get("/profile", middlewares.RequireAuthMiddleware, controllers.Profile)
	userApi.Put("/update-address", middlewares.RequireAuthMiddleware, controllers.UpdateAddress)

	productApi := api.Group("/products")
	productApi.Post("/", middlewares.RequireAuthMiddleware, controllers.CreateProduct)
}