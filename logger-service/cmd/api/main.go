package main

import (
	"context"
	"log"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "5001"
)

var client *mongo.Client

type Config struct{}

func main() {
	// connect to MongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.panic(err)
	}
	client = mongoClient
}

func connectToMongo() (*mongo.Client, error) {
	// create the connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	return c, nil
}
