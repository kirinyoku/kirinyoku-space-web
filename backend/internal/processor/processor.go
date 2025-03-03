// Package processor provides functionality for transforming raw messages into structured data.
// It handles the parsing, validation, and conversion of unstructured message content
// into well-defined data structures that can be easily stored.
package processor

import (
	"fmt"
	"log"
	"strings"

	"github.com/kirinyoku/kirinyoku-space-web/backend/internal/bot"
)

// ProcessedMessage represents a fully processed message ready for storage or further handling.
// It contains structured data extracted from the original message text.
type ProcessedMessage struct {
	Name string   // The name or title of the resource
	Type string   // The type or category of the resource
	Tags []string // List of tags associated with the resource
	URL  string   // URL linking to the resource
}

// Processor handles the transformation of raw bot messages into structured data.
// It operates asynchronously using channels for input and output communication.
type Processor struct {
	inputChan  chan bot.Message      // Channel for receiving raw messages
	outputChan chan ProcessedMessage // Channel for sending processed messages
}

// NewProcessor creates and initializes a new Processor with the specified input and output channels.
// Returns an error if either channel is nil.
func NewProcessor(inputChan chan bot.Message, outputChan chan ProcessedMessage) (*Processor, error) {
	if inputChan == nil || outputChan == nil {
		return nil, fmt.Errorf("input and output channels cannot be nil")
	}

	return &Processor{
		inputChan:  inputChan,
		outputChan: outputChan,
	}, nil
}

// Start begins the message processing loop in a separate goroutine.
// It continuously reads from the input channel, processes each message,
// and sends the results to the output channel.
func (p *Processor) Start() {
	go func() {
		for msg := range p.inputChan {
			processed, err := p.processMessage(msg)
			if err != nil {
				log.Printf("Skipping message due to processing error: %v", err)
				continue
			}

			select {
			case p.outputChan <- *processed:
				log.Printf("Processed message: %+v", processed)
			default:
				log.Printf("Output channel full, skipping message: %s", processed.Name)
			}
		}
	}()

	log.Printf("Processor started")
}

// processMessage transforms a raw bot message into a structured ProcessedMessage.
// It parses the message text line by line, extracting key-value pairs and validating
// that all required fields are present. It also ensures the message contains a valid URL.
func (p *Processor) processMessage(msg bot.Message) (*ProcessedMessage, error) {
	lines := strings.Split(msg.Text, "\n")
	fields := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		if key != "" && value != "" {
			fields[key] = value
		}
	}

	requiredFields := []string{"name", "type", "tags"}
	for _, field := range requiredFields {
		if value, ok := fields[field]; !ok || strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("missing or empty required field: %s", field)
		}
	}

	// Use the URL from the bot's Message struct
	if msg.URL == "" {
		return nil, fmt.Errorf("no valid URL found in message")
	}

	tags := parseTags(fields["tags"])
	if len(tags) == 0 {
		return nil, fmt.Errorf("no valid tags found after parsing")
	}

	return &ProcessedMessage{
		Name: fields["name"],
		Type: fields["type"],
		Tags: tags,
		URL:  msg.URL,
	}, nil
}

// parseTags converts a raw tags string into a clean slice of individual tags.
// It handles various separator formats (spaces, commas, etc.), removes any
// leading '#' symbols, and filters out empty tags. Returns nil if the input
// is empty or no valid tags are found.
func parseTags(rawTags string) []string {
	if rawTags == "" {
		return nil
	}

	// Split by spaces, commas, or other common separators
	separators := []string{" ", ",", ";", "|"}
	var tags []string

	current := rawTags
	for _, sep := range separators {
		if strings.Contains(current, sep) {
			tags = strings.Split(current, sep)
			break
		}
	}

	// If no separators found, treat as single tag
	if tags == nil {
		tags = []string{rawTags}
	}

	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		// Remove # prefix if present
		tag = strings.TrimPrefix(tag, "#")
		tag = strings.TrimSpace(tag)

		if tag != "" {
			result = append(result, tag)
		}
	}

	return result
}
