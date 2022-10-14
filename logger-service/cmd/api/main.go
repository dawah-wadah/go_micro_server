package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = 5001
	mongoURL = "mongodb://localhost:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	// 5 - Set up config with new models
	Models data.Models
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

	app := Config{
		Models: data.New(client),
	}

	// start web server

	// go app.serve()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.panic(err)
	}

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
