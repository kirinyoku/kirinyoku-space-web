package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	Text   string
	ChatID int64
}

type Bot struct {
	api         *tgbotapi.BotAPI
	channelID   int64
	messageChan chan Message
}

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

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message.Chat.ID != b.channelID {
				continue
			}

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
