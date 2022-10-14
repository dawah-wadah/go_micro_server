package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct{}

// define kind of payload youre going to receive
type RPCPayload struct {
	Name string
	Data string
}

// write methods that we want to expose
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(
		context.TODO(),
		data.LogEntry{
			Name:      payload.Name,
			Data:      payload.Data,
			CreatedAt: time.Now(),
		})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}
	*resp = "Processed payload via RPC: " + payload.Name
	return nil
}

// ## Steps to implement RPC in Go
// - In the Microservice that will receive the call, implement a server
// - declare a type thats going to be the RPC Server
// - then we declare a type that will be the kind of data we're going to receive for any methods that are tied to rpc server
// 		- then we declare the function (in our current case, just one function) that has the reciever pointer (&RPCServer)
// 		- and that will require a payload of some sort, and resp is a pointer to a string that will send a message back to the people who call us
//		- if there is no error, we use the response of the database call and set the pointer to the response string to equal the value of our return message + some response specific attribute (ie `Name`)
