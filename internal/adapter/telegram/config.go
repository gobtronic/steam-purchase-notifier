package telegram

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	BotToken string
	ChatID   int64
}

func LoadConfig() (*Config, error) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Please set the TELEGRAM_BOT_TOKEN environment variable with your Telegram bot token")
	}
	envTelegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	if envTelegramChatID == "" {
		log.Fatal("Please set the TELEGRAM_CHAT_ID environment variable with your Telegram chat ID")
	}
	telegramChatID, err := strconv.Atoi(envTelegramChatID)
	if err != nil {
		log.Fatal(err)
	}
	chatID := int64(telegramChatID)
	return &Config{
		BotToken: botToken,
		ChatID:   chatID,
	}, nil
}
