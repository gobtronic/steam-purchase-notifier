package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"time"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
)

type DiscordNotifier struct {
	botToken  string
	channelID string
}

func NewDiscordNotifier() (*DiscordNotifier, error) {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Please set the DISCORD_BOT_TOKEN environment variable with your Discord bot token")
	}

	envDiscordChannelID := os.Getenv("DISCORD_CHANNEL_ID")
	if envDiscordChannelID == "" {
		log.Fatal("Please set the DISCORD_CHANNEL_ID environment variable with your Discord channel ID")
	}

	return &DiscordNotifier{
		botToken:  botToken,
		channelID: envDiscordChannelID,
	}, nil
}

func (n *DiscordNotifier) Notify(game domain.Game) error {

	text := fmt.Sprintf(
		"💸 A new game is available in your Steam Family library!\n\n**[%s](%s)**",
		game.Name,
		game.StoreURL,
	)

	return n.notify(text, dictionary{})
}

type dictionary = map[string]any

func (n *DiscordNotifier) notify(content string, additionalInfo dictionary) error {
	payload, err := n.buildPayload(content, additionalInfo)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(payload)
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", n.channelID)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bot "+n.botToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord API returned %s: %s", resp.Status, string(respBody))
	}

	return nil
}

func (n *DiscordNotifier) buildPayload(content string, additionalInfo dictionary) ([]byte, error) {
	data := dictionary{
		"content": content,
	}

	maps.Copy(data, additionalInfo)

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}
