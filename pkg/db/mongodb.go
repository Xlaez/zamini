package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectWithMongo() *mongo.Client{
	if err := godotenv.Load(); err !=nil{
		log.Fatal(err)
	}

	uri := os.Getenv("MONGO_URI")

	if uri == ""{
		log.Fatal("MONGO_URI is not a valid environmental variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err !=nil{
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err !=nil{
		log.Fatal(err)
	}else{
		log.Println("Connected to MongoDB")
	}

	return client
}

var Client *mongo.Client = ConnectWithMongo()

func OpenCollection(client *mongo.Client, collectionName string)*mongo.Collection{
	var collection *mongo.Collection = client.Database("ecommerce").Collection(collectionName)
	return collection
}
