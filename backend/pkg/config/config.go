// Package config provides configuration management functionality for the application.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration parameters for the application.
type Config struct {
	TelegramToken   string
	TelegramChatID  int64
	MongoURI        string
	MongoDatabase   string
	MongoCollection string
	APIPort         string
}

// Load reads configuration from environment variables and returns a Config struct.
// It returns an error if required configuration is missing or invalid.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
		TelegramChatID:  parseChatID(os.Getenv("TELEGRAM_CHAT_ID")),
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDatabase:   os.Getenv("MONGO_DATABASE"),
		MongoCollection: os.Getenv("MONGO_COLLECTION"),
		APIPort:         os.Getenv("API_PORT"),
	}

	if cfg.APIPort == "" {
		cfg.APIPort = ":8080"
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// parseChatID converts a string chat ID to int64.
// Returns 0 if the input is empty or invalid.
func parseChatID(chatIDStr string) int64 {
	if chatIDStr == "" {
		return 0
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0
	}

	return chatID
}

// validate checks if all required configuration fields are properly set.
// Returns an error if any required field is missing or invalid.
func (c *Config) validate() error {
	if c.TelegramToken == "" {
		return fmt.Errorf("TELEGRAM_TOKEN is required")
	}

	if c.TelegramChatID == 0 {
		return fmt.Errorf("TELEGRAM_CHAT_ID is required and must be a valid integer")
	}

	if c.MongoURI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}

	if c.MongoDatabase == "" {
		return fmt.Errorf("MONGO_DATABASE is required")
	}

	if c.MongoCollection == "" {
		return fmt.Errorf("MONGO_COLLECTION is required")
	}

	return nil
}
