package main

import (
	"context"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = 5001
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
}

func main() {
	// Install mongo drivers
	// connect to Mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.panic(err)
	}
	client = mongoClient

	// create a context inorder to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.panic(err)
		}
	}()
}

func connectToMongo() (*mongo.Client, error) {
	// create the connection options

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credentials{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.PrintLn("Error Conencting:", err)
	}

	return client, err
}
