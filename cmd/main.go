package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/gamestore"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/steam"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/telegram"
	"github.com/gobtronic/steam-purchase-notifier/internal/usecase"
	"github.com/joho/godotenv"
)

type config struct {
	steamAPIKey      string
	steamID          string
	telegramBotToken string
	telegramChatID   int64
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	steamClient := steam.NewSteamClient(cfg.steamAPIKey, cfg.steamID, http.DefaultClient)
	gameStoreFilePath := path.Join(os.Getenv("GOPATH"), "gamelist.json")
	gameStore := gamestore.NewGameStore(gameStoreFilePath)
	notifier := telegram.NewTelegramNotifier(cfg.telegramBotToken, cfg.telegramChatID)

	games, err := steamClient.FetchGames()
	if err != nil {
		log.Fatal(err)
	}
	newGames, err := usecase.FilterNewGames(games, gameStore)
	if err != nil {
		log.Print(err)
	}
	gameStore.Write(games)

	if len(games) != len(newGames) {
		usecase.NotifyGames(newGames, notifier)
	}
}

func loadConfig() (config, error) {
	godotenv.Load()
	steamAPIKey := os.Getenv("STEAM_API_KEY")
	if steamAPIKey == "" {
		log.Fatal("Please set the STEAM_API_KEY environment variable with your Steam API key")
	}
	steamID := os.Getenv("STEAM_ID")
	if steamID == "" {
		log.Fatal("Please set the STEAM_ID environment variable with Steam ID that will be watched")
	}
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
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

	return config{
		steamAPIKey:      steamAPIKey,
		steamID:          steamID,
		telegramBotToken: telegramBotToken,
		telegramChatID:   int64(telegramChatID),
	}, nil
}
