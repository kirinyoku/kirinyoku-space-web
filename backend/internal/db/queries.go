package db

import (
	"context"
	"strings"
	"time"

	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/processor"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// PostsResponse represents the response structure for post queries
type PostsResponse struct {
	Posts      []processor.ProcessedMessage
	TotalCount int64
}

// GetPostsWithFilters retrieves posts with specified filters and pagination
// query: search term for post name
// tag: specific tag to filter by
// postType: type of post to filter
// language: language tag to filter
// page: page number for pagination
// limit: number of posts per page
func (d *DB) GetPostsWithFilters(query, tag, postType, language string, page, limit int) (PostsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	filter := bson.M{}
	if query != "" {
		filter["name"] = bson.M{"$regex": query, "$options": "i"}
	}
	if postType != "" {
		filter["type"] = postType
	}

	var tagConditions []bson.M
	if tag != "" {
		tagConditions = append(tagConditions, bson.M{"tags": tag})
	}
	if language != "" {
		tagConditions = append(tagConditions, bson.M{"tags": bson.M{"$regex": language + "$"}})
	}
	if len(tagConditions) > 0 {
		if len(tagConditions) == 1 {
			filter["tags"] = tagConditions[0]["tags"]
		} else {
			filter["$and"] = tagConditions
		}
	}

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return PostsResponse{}, err
	}

	cur, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return PostsResponse{}, err
	}
	defer cur.Close(ctx)

	var posts []processor.ProcessedMessage
	for cur.Next(ctx) {
		var msg processor.ProcessedMessage
		if err := cur.Decode(&msg); err != nil {
			return PostsResponse{}, err
		}
		for i, tag := range msg.Tags {
			msg.Tags[i] = strings.TrimPrefix(tag, "#")
		}
		posts = append(posts, msg)
	}

	if err := cur.Err(); err != nil {
		return PostsResponse{}, err
	}

	return PostsResponse{Posts: posts, TotalCount: total}, nil
}

// GetPosts retrieves all posts with pagination
func (d *DB) GetPosts(page, limit int) (PostsResponse, error) {
	return d.GetPostsWithFilters("", "", "", "", page, limit)
}

// GetPostsByTag retrieves posts with a specific tag
func (d *DB) GetPostsByTag(tag string, page, limit int) (PostsResponse, error) {
	return d.GetPostsWithFilters("", tag, "", "", page, limit)
}

// GetPostsByType retrieves posts of a specific type
func (d *DB) GetPostsByType(postType string, page, limit int) (PostsResponse, error) {
	return d.GetPostsWithFilters("", "", postType, "", page, limit)
}

// GetPostsByLanguage retrieves posts in a specific language
func (d *DB) GetPostsByLanguage(language string, page, limit int) (PostsResponse, error) {
	return d.GetPostsWithFilters("", "", "", language, page, limit)
}

// SearchPosts searches posts by query string
func (d *DB) SearchPosts(query string, page, limit int) (PostsResponse, error) {
	return d.GetPostsWithFilters(query, "", "", "", page, limit)
}

// GetTags retrieves all unique tags from the collection
func (d *DB) GetTags() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{"$unwind", "$tags"}},
		{{"$group", bson.D{{"_id", "$tags"}}}},
		{{"$sort", bson.D{{"_id", 1}}}},
		{{"$project", bson.D{{"tag", "$_id"}, {"_id", 0}}}},
	}

	cur, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tags []string
	for cur.Next(ctx) {
		var result struct {
			Tag string `bson:"tag"`
		}
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		tags = append(tags, strings.TrimPrefix(result.Tag, "#"))
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if tags == nil {
		return []string{}, nil
	}
	return tags, nil
}

// GetLanguages retrieves all unique language codes from tags
func (d *DB) GetLanguages() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{"$project", bson.D{{"lastTag", bson.M{"$arrayElemAt": []interface{}{"$tags", -1}}}}}},
		{{"$match", bson.D{{"lastTag", bson.M{"$regex": "^[a-z]{2}$"}}}}},
		{{"$group", bson.D{{"_id", "$lastTag"}}}},
		{{"$sort", bson.D{{"_id", 1}}}},
		{{"$project", bson.D{{"language", "$_id"}, {"_id", 0}}}},
	}

	cur, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var languages []string
	for cur.Next(ctx) {
		var result struct {
			Language string `bson:"language"`
		}
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		languages = append(languages, result.Language)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if languages == nil {
		return []string{}, nil
	}
	return languages, nil
}
