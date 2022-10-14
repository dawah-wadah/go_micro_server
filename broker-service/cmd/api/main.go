package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	// drivers to cennect to rabbitmq
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	RabbitMQ *amqp.Connection
}

func main() {
	// connect to rabbitmq
	// try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// install a third party routing package
	app := Config{
		RabbitMQ: rabbitConn,
	}

	log.Printf("Starting Broker Service on port %s\n", webPort)

	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	// Note: RabbitMQ is slow to start but once it starts its rock solid
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbit is ready

	for {

		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready: %d", counts)
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing Off")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
