// Package bot provides functionality for interacting with the Telegram Bot API
package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Message represents a message received from Telegram
type Message struct {
	Text   string // The text content of the message
	ChatID int64  // The ID of the chat the message was sent in
}

// Bot handles the Telegram bot functionality
type Bot struct {
	api         *tgbotapi.BotAPI // The Telegram Bot API client
	channelID   int64            // ID of the channel to listen to
	messageChan chan Message     // Channel for sending received messages
}

// New creates a new Bot instance with the given token and configuration
func New(token string, channelID int64, messageChan chan Message) *Bot {
	botapi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot API: %v", err)
	}

	return &Bot{
		api:         botapi,
		channelID:   channelID,
		messageChan: messageChan,
	}
}

// Start begins listening for messages from the configured Telegram channel
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			// Skip messages not from our target channel
			if update.Message.Chat.ID != b.channelID {
				continue
			}

			// Skip nil messages
			if update.Message == nil {
				continue
			}

			msg := Message{
				Text:   update.Message.Text,
				ChatID: update.Message.Chat.ID,
			}

			select {
			case b.messageChan <- msg:
				log.Printf("Message sent to channel: %s", msg.Text)
			default:
				log.Printf("Message channel full, skipping message: %s", msg.Text)
			}
		}
	}()

	log.Printf("Bot started, listening for messages from channel %d", b.channelID)
}
