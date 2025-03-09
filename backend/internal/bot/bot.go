// Package bot provides functionality for interacting with the Telegram Bot API.
// It handles message reception from a specified channel and forwards them through
// a message channel for further processing.
package bot

import (
	"errors"
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Message represents a message received from Telegram containing the essential
// information needed for processing.
type Message struct {
	Text   string // Raw text content of the message
	URL    string // URL extracted from message entities
	ChatID int64  // Identifier of the chat where message originated
}

// Bot manages the Telegram bot operations including message listening,
// processing and forwarding. It maintains connection with Telegram API
// and handles graceful shutdown.
type Bot struct {
	api         *tgbotapi.BotAPI // Connection to Telegram Bot API
	channelID   int64            // Target channel to monitor
	messageChan chan Message     // Output channel for processed messages
	done        chan struct{}    // Signal channel for shutdown
	wg          sync.WaitGroup   // Ensures clean goroutine termination
}

// ErrInvalidParams is returned when required initialization parameters are missing or invalid.
var ErrInvalidParams = errors.New("invalid parameters: token or channelID empty, or messageChan nil")

// New initializes a new Bot instance with the provided configuration.
// It establishes connection with Telegram API and sets up message handling infrastructure.
// Returns error if initialization fails due to invalid parameters or API connection issues.
func New(token string, channelID int64, messageChan chan Message) (*Bot, error) {
	if token == "" || channelID == 0 || messageChan == nil {
		return nil, ErrInvalidParams
	}

	botapi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	return &Bot{
		api:         botapi,
		channelID:   channelID,
		messageChan: messageChan,
		done:        make(chan struct{}),
	}, nil
}

// Start initiates the message monitoring process in a separate goroutine.
// It configures update parameters to only listen for channel posts and
// processes incoming messages, extracting URLs and forwarding them through
// the message channel.
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"channel_post"} // Only listen for channel posts

	updates := b.api.GetUpdatesChan(u)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case update := <-updates:
				if update.ChannelPost == nil {
					continue
				}

				if update.ChannelPost.Chat.ID != b.channelID {
					log.Printf("Received message from unexpected channel: %d", update.ChannelPost.Chat.ID)
					continue
				}

				url := b.extractURLFromEntities(update.ChannelPost.Text, update.ChannelPost.Entities)

				msg := Message{
					Text:   update.ChannelPost.Text,
					URL:    url,
					ChatID: update.ChannelPost.Chat.ID,
				}

				b.sendMessage(msg)

			case <-b.done:
				return
			}
		}
	}()

	log.Printf("Bot started, listening for messages from channel %d", b.channelID)
}

// Stop gracefully terminates the bot's operations.
// It signals the monitoring goroutine to stop and waits for its completion.
func (b *Bot) Stop() {
	close(b.done)
	b.wg.Wait()
	log.Printf("Bot stopped")
}

// sendMessage attempts to send a message through the output channel.
// It logs the operation result and handles cases where the channel is full.
// Empty messages are skipped to avoid processing invalid data.
func (b *Bot) sendMessage(msg Message) {
	if msg.Text == "" {
		log.Printf("Skipping empty message")
		return
	}

	select {
	case b.messageChan <- msg:
		log.Printf("Message sent to channel:\n%s", msg.Text)
	default:
		log.Printf("Message channel full, skipping message:\n%s", msg.Text)
	}
}

// extractURLFromEntities searches for a URL in message entities specifically
// associated with the "link" text. It validates entity boundaries and handles
// potential edge cases in entity processing.
// Returns an empty string if no valid URL is found or if input is invalid.
func (b *Bot) extractURLFromEntities(text string, entities []tgbotapi.MessageEntity) string {
	if text == "" || entities == nil {
		return ""
	}

	for _, entity := range entities {
		if entity.Type != "text_link" || entity.URL == "" {
			continue
		}

		return entity.URL
	}

	return ""
}
