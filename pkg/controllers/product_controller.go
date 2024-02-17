package controllers

import (
	"context"
	"ecommerce/pkg/db"
	"ecommerce/pkg/models"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var productionCollection = db.OpenCollection(db.Client,"products")

func CreateProduct(c *fiber.Ctx) error{
	gofakeit.Seed(0)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	userType, err := c.Locals("userType").(string)

	if !err{
			return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data for model binding",
			"data":    err,
		})
	}

	if userType != "ADMIN"{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Only an Admin can add a product",
		})
	}

	type RequestBody struct{
		Category string `json:"category" validate:"required"`
		Name string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
		Price float64 `json:"price" validate:"required"`
		AvailabeQuantity int16 `json:"quantity" validate:"required"`
	}

	var requestBody RequestBody

	if err := c.BodyParser(&requestBody);err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "cannot bind request body",
			"data": err.Error(),
		})
	}

	var product models.Product
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.Category = requestBody.Category
	product.AvailableQuantity = int(requestBody.AvailabeQuantity)
	product.Description = requestBody.Description
	product.Images =  make([]string, 0)

	filter := bson.D{{Key: "name", Value: product.Name}}
	if _, err := productionCollection.FindOne(ctx, filter).DecodeBytes(); err ==nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "cannot create two products with the same name",
		})
	}  

	if _, err := productionCollection.InsertOne(ctx, product); err !=nil{
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "cannot add product to databaser",
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status": "success",
		"message": fmt.Sprintf(product.Name, "added to the database"),
		"data": product,
	})

}