package main

import (
	"context"
	"ecommerce/pkg/db"
	"ecommerce/pkg/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	port = ":9100"
)

func main(){
	app := fiber.New()
	app.Use(logger.New())
	routes.Router(app)

	client := db.Client

	defer func(client *mongo.Client, ctx context.Context){
		err := client.Disconnect(ctx)
		if err !=nil{
			log.Fatal(err)
		}
	}(client, context.Background())

	if err := app.Listen(port); err !=nil{
		log.Fatal("Error starting the server: ", err.Error())
	}
}