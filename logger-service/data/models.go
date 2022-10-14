package data

import (
	"context"
	"time"
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
		log.PrintLn("Error inserting log entry into logs: ", err)
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
	if err!= nil {
        log.PrintLn("Finding all docs error:", err)
		return nil, err
    }

	defer cursor.Close()

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err = cursor.Decode(&item)
        if err!= nil {
			log.PrintLn("Error decoding log into slice:", err)
            return nil, err
	} else {
		logs = append(logs, &item)
	}

	return logs, nil
}

// get record by id
// delete collection
