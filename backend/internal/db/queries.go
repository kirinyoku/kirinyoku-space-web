package db

import (
	"context"
	"fmt"
	"time"

	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/processor"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GetPots retrieves a paginated list of posts from the database.
// It accepts page number and limit parameters to implement pagination.
// Returns a slice of ProcessedMessage and an error if any occurs.
func (db *DB) GetPosts(page int, limit int) ([]processor.ProcessedMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cur, err := db.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var posts []processor.ProcessedMessage
	for cur.Next(ctx) {
		var message processor.ProcessedMessage
		if err := cur.Decode(&message); err != nil {
			return nil, err
		}

		posts = append(posts, message)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// SearchPosts performs a text search on the database using the provided query string.
// It supports pagination through page number and limit parameters.
// Returns matching posts as a slice of ProcessedMessage and an error if any occurs.
func (db *DB) SearchPosts(query string, page, limit int) ([]processor.ProcessedMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * limit

	filter := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: query}}}}
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cur, err := db.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var posts []processor.ProcessedMessage
	for cur.Next(ctx) {
		var message processor.ProcessedMessage
		if err := cur.Decode(&message); err != nil {
			return nil, err
		}

		posts = append(posts, message)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostsByTag retrieves posts that contain a specific tag.
// It supports pagination through page number and limit parameters.
// Returns matching posts as a slice of ProcessedMessage and an error if any occurs.
func (db *DB) GetPostsByTag(tag string, page, limit int) ([]processor.ProcessedMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * limit

	tag = fmt.Sprintf("#%s", tag)

	filter := bson.D{{Key: "tags", Value: bson.D{{Key: "$eq", Value: tag}}}}
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cur, err := db.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var posts []processor.ProcessedMessage
	for cur.Next(ctx) {
		var message processor.ProcessedMessage
		if err := cur.Decode(&message); err != nil {
			return nil, err
		}

		posts = append(posts, message)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetTags retrieves all unique tags in the collection.
// Returns a slice of strings containing all unique tags and an error if any occurs.
func (db *DB) GetTags() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.collection.Distinct(ctx, "tags", bson.D{})

	var tags []interface{}
	if err := result.Decode(&tags); err != nil {
		return nil, err
	}

	tagList := make([]string, len(tags))
	for i, tag := range tags {
		if tagStr, ok := tag.(string); ok {
			tagList[i] = tagStr
		} else {
			continue
		}
	}

	return tagList, nil
}
