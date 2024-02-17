package controllers

import (
	"context"
	"ecommerce/pkg/db"
	"ecommerce/pkg/helpers"
	"ecommerce/pkg/models"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userController *mongo.Collection = db.OpenCollection(db.Client, "users")

func Signup(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	var user models.User
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	var userAddress models.Address
	userAddress.ZipCode = gofakeit.Zip()
	userAddress.City = gofakeit.City()
	userAddress.State = gofakeit.State()
	userAddress.Country = gofakeit.Country()
	userAddress.HouseNumber = gofakeit.StreetNumber()
	userAddress.Street = gofakeit.Street()
	user.Address = userAddress
	user.Orders = make([]models.Order, 0)
	user.UserCart = make([]models.ProductsToOrder, 0)
	
	if err := c.BodyParser(&user); err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid user data request",
			"data": err.Error(),
		})
	}

	if user.Password == os.Getenv("ADMIN_PASS") && user.Email == os.Getenv("ADMIN_EMAIL"){
		user.UserType = "ADMIN"
	}else{
		user.UserType = "USER"
	}

	if user.UserType == "ADMIN"{
		filter := bson.M{"userType": "ADMIN"}
		if _, err := userController.FindOne(ctx, filter).DecodeBytes(); err !=nil{
			return c.Status(400).JSON(fiber.Map{
				"status": "error",
				"message": "Admin already exists",
				"data": c.JSON(err),
			})
		}
	}

	filter := bson.M{"email": user.Email}
	if _, err := userController.FindOne(ctx, filter).DecodeBytes(); err ==  nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "User with email already registerd",
			"data": c.JSON(nil),
		})
	}

	password, err := helpers.HashPassword(user.Password)

	if err !=nil{
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Cannot hash password",
			"data": c.JSON(nil),
		})
	}

	user.Password = password

	if _, err := userController.InsertOne(ctx, user); err !=nil{
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "cannot create user",
			"data": err.Error(),
		})
	}

	signedToken, err := helpers.CreateToken(user.ID, user.UserType)

	if err !=nil{
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "cannot create auth tokens",
			"data": err.Error(),
		})
	}


	cookie := &fiber.Cookie{
		Name: "x-auth-jwt",
		Value: signedToken,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	}

	c.Cookie(cookie)

	return c.Status(200).JSON(fiber.Map{
		"status": "succees",
		"message": "User created successfully",
		"data": helpers.SterializeUser(user),
	})
}

func SingIn(c *fiber.Ctx)error{
	type SignInRequest struct{
		Email string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,alphanum"`
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var req SignInRequest

	if err := c.BodyParser(&req); err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid request body",
			"data": err.Error(),
		})
	}

	var existingUser bson.Raw
	
	if err := userController.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser);err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "user does not exist",
			"data": err.Error(),
		})
	}

	isValid := helpers.ComparePassword(req.Password, existingUser.Lookup("password").StringValue())

	if isValid != nil{
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid credentials",
			"data": nil,
		})
	}

	signedToken, err := helpers.CreateToken(existingUser.Lookup("_id").ObjectID(), existingUser.Lookup("userType").StringValue())

	if err !=nil{
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to create auth token",
			"data": err.Error(),
		})
	}

	cookie := &fiber.Cookie{
		Name: "x-auth-jwt",
		Value: signedToken,
		Expires: time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(cookie)

	user := models.User{}

	user.ID = existingUser.Lookup("_id").ObjectID()
    user.Username = existingUser.Lookup("username").StringValue()
    user.Email = existingUser.Lookup("email").StringValue()
    user.UserType = existingUser.Lookup("userType").StringValue()
    user.CreatedAt = existingUser.Lookup("createdAt").Time()
    user.UpdatedAt = existingUser.Lookup("updatedAt").Time()
	
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User signin successfully",
		"data": user,
	})
}