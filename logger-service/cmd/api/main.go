package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = 5001
	mongoURL = "mongodb://mongo:27017"
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
		log.Panic(err)
	}
	client = mongoClient

	// create a context inorder to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}
	// we have to register the rpc server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()
	// start web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) rpcListen() error {
	log.Printf("Starting rpc server on port %s", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			// in this case, we just want to start over
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connectToMongo() (*mongo.Client, error) {
	// create the connection options

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error Conencting:", err)
	}

	return client, err
}
