// Package db provides functionality for storing processed messages in MongoDB.
// It handles database connections, data persistence, index management and CRUD operations.
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

// DB manages MongoDB operations including connection management,
// data storage, and asynchronous message processing.
type DB struct {
	client     *mongo.Client                   // MongoDB client connection
	collection *mongo.Collection               // Target collection for storing messages
	inputChan  chan processor.ProcessedMessage // Channel for receiving processed messages
}

// New creates and initializes a new DB instance with the specified MongoDB connection parameters.
// It establishes a connection to MongoDB, verifies connectivity with a ping test,
// and sets up the required collection and indexes.
// Returns an error if connection, ping, or index creation fails.
func New(uri, database, collection string, inputChan chan processor.ProcessedMessage) (*DB, error) {
	// Configure client options with connection timeout
	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)

	// Establish connection to MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection with ping test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// Get reference to the specified collection
	coll := client.Database(database).Collection(collection)

	// Initialize DB instance
	db := &DB{
		client:     client,
		collection: coll,
		inputChan:  inputChan,
	}

	// Create necessary indexes for efficient querying
	if err := db.createIndex(); err != nil {
		return nil, err
	}

	return db, nil
}

// Start begins the message processing loop in a separate goroutine.
// It continuously reads from the input channel and persists each message to MongoDB.
// Errors during save operations are logged but don't interrupt processing.
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

// saveMessage persists a processed message to MongoDB.
// It converts the message to BSON format and inserts it into the collection.
// Uses a timeout context to prevent hanging operations.
func (db *DB) saveMessage(message processor.ProcessedMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert message to BSON document format
	doc := bson.D{
		{Key: "name", Value: message.Name},
		{Key: "type", Value: message.Type},
		{Key: "tags", Value: message.Tags},
		{Key: "url", Value: message.URL},
	}

	// Insert document into collection
	_, err := db.collection.InsertOne(ctx, doc)
	return err
}

// createIndex sets up MongoDB indexes to optimize query performance.
// Creates a multi-key index on tags for efficient tag-based lookups
// and a text index on name for text search capabilities.
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

	// Create both indexes in a single operation
	_, err := db.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{tagsIndex, nameIndex})
	if err != nil {
		return err
	}

	log.Println("Indexes created on tags and name")
	return nil
}

// Disconnect cleanly closes the MongoDB connection.
// Should be called when the application is shutting down to release resources.
func (db *DB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.client.Disconnect(ctx)
}
