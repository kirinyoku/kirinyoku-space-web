package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/api"
	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/bot"
	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/db"
	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/processor"
	"github.com/kirinyoku/kirinyoku-space-web/backend/pkg/config"
)

// main initializes and starts all application components in the following order:
// 1. Load configuration from environment variables
// 2. Create communication channels between components
// 3. Initialize and start the Telegram bot
// 4. Initialize and start the message processor
// 5. Initialize and start the MongoDB connection
// 6. Start the HTTP API server
// 7. Wait for shutdown signal
func main() {
	// Load application configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create buffered channels for inter-component communication
	botChan := make(chan bot.Message, 100)                 // Channel for raw messages from Telegram
	procChan := make(chan processor.ProcessedMessage, 100) // Channel for processed messages

	// Initialize and start Telegram bot
	bot, err := bot.New(cfg.TelegramToken, cfg.TelegramChatID, botChan)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	bot.Start()

	// Initialize and start message processor
	processor, err := processor.NewProcessor(botChan, procChan)
	if err != nil {
		log.Fatalf("Failed to create processor: %v", err)
	}
	processor.Start()

	// Initialize and start MongoDB connection
	db, err := db.New(cfg.MongoURI, cfg.MongoDatabase, cfg.MongoCollection, procChan)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Disconnect()
	db.Start()

	// Initialize and start HTTP API server
	server := api.NewServer(db)
	go func() {
		if err := server.Start(cfg.APIPort); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// Set up graceful shutdown on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
