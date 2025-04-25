package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client    *mongo.Client
	database  *mongo.Database
	surveys   *mongo.Collection
	questions *mongo.Collection
	responses *mongo.Collection
)

// InitDB initializes MongoDB connection
func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get Atlas connection string from environment variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	uri := os.Getenv("MONGODB_URI")

	// Set client options and connect
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB Atlas: %v", err)
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB Atlas: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB Atlas!")

	// Set global variables
	Client = client
	database = Client.Database("OSP_backend")
	surveys = database.Collection("surveys")
	questions = database.Collection("questions")
	responses = database.Collection("responses")
}

// GetCollection retrieves a MongoDB collection by name
func GetCollection(name string) *mongo.Collection {
	return database.Collection(name)
}
