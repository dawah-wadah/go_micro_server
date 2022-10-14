package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}

}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// this is the `receiver` syntax, meaning this function will make use of this context
func (l *LogEntry) Insert(entry LogEntry) (error error) {
	// declare a collection variable. can think of them as tables
	collection := client.Database("logs").Collection("logs")

	// only differnce is we use insert one
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry into logs: ", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err = cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}

	}
	return logs, nil
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*15)
}

// get record by id
func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := getContext()
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	// convert id to something mongo can injest
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// delete collection

func (l *LogEntry) DropCollection() error {
	ctx, cancel := getContext()
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping collection:", err)
		return err
	}
	return nil
}

// this is actually my own blind attemp at updating, it is wrong but i will keep it to see how i think
// func (l *LogEntry) UpdateOne(entry LogEntry) (LogEntry, error) {

// 	collection := client.Database("logs").Collection("logs")
// 	err := collection.UpdateOne(ctx, LogEntry{
// 		ID:        entry.ID,
// 		Data:      entry.Data,
// 		CreatedAt: entry.CreatedAt,
// 		UpdatedAt: time.Now(),
// 	})
// 	if err != nil {
// 		log.Println("Error updating log entry into logs: ", err)
// 		return LogEntry{}, err
// 	}

// }

// this is the official way to update a log entry
func (l *LogEntry) UpdateOne() (*mongo.UpdateResult, error) {
	ctx, cancel := getContext()
	defer cancel()
	collection := client.Database("logs").Collection("logs")

	// get the docID from the reciever
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}
	return result, nil
}
