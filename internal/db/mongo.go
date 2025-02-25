package db

import (
	"context"
	"log"
	"time"

	"github.com/kirinyoku/kirinyoku-space-tg/internal/processor"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type DB struct {
	client     *mongo.Client
	collection *mongo.Collection
	inputChan  chan processor.ProcessedMessage
}

func New(uri, database, collection string, inputChan chan processor.ProcessedMessage) (*DB, error) {
	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	coll := client.Database(database).Collection(collection)

	db := &DB{
		client:     client,
		collection: coll,
		inputChan:  inputChan,
	}

	if err := db.createIndex(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Start() {
	go func() {
		for message := range db.inputChan {
			if err := db.saveMessage(message); err != nil {
				log.Printf("Failed to save message %s: %v", message.Name, err)
				continue
			}

			log.Printf("Saved message\n%s", message)
		}
	}()

	log.Println("DB started, listening for processed messages")
}

func (db *DB) saveMessage(message processor.ProcessedMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.D{
		{Key: "name", Value: message.Name},
		{Key: "type", Value: message.Type},
		{Key: "tags", Value: message.Tags},
		{Key: "url", Value: message.URL},
	}

	_, err := db.collection.InsertOne(ctx, doc)
	return err
}

func (db *DB) createIndex() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Index on tags (multi-key index for arrays)
	tagsIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "tags", Value: 1}},
	}
	// Index on name (text index for search)
	nameIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}},
	}

	_, err := db.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{tagsIndex, nameIndex})
	if err != nil {
		return err
	}

	log.Println("Indexes created on tags and name")
	return nil
}

func (db *DB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.client.Disconnect(ctx)
}
