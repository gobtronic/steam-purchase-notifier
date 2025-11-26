package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/gamestore"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/steam"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/telegram"
	"github.com/gobtronic/steam-purchase-notifier/internal/port"
	"github.com/gobtronic/steam-purchase-notifier/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var telegramNotifier bool
var rootCmd = &cobra.Command{
	Use:   "steam-purchase-notifier",
	Short: "Watch a Steam account purchases through notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !telegramNotifier {
			return fmt.Errorf("you must specify at least one notifier flag, e.g. --telegram")
		}

		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		var notifiers []port.Notifier
		if telegramNotifier {
			notifier := telegram.NewTelegramNotifier(*cfg.telegramBotToken, *cfg.telegramChatID)
			notifiers = append(notifiers, notifier)
		}

		steamClient := steam.NewSteamClient(cfg.steamAPIKey, cfg.steamID, http.DefaultClient)
		gameStoreFilePath := path.Join(os.Getenv("GOPATH"), "gamelist.json")
		gameStore := gamestore.NewGameStore(gameStoreFilePath)
		games, err := steamClient.FetchGames()
		if err != nil {
			return err
		}
		newGames, err := usecase.FilterNewGames(games, gameStore)
		if err != nil {
			return err
		}
		gameStore.Write(games)

		if len(games) != len(newGames) {
			for _, notifier := range notifiers {
				usecase.NotifyGames(newGames, notifier)
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&telegramNotifier, "telegram", false, "telegram")
}

type config struct {
	steamAPIKey      string
	steamID          string
	telegramBotToken *string
	telegramChatID   *int64
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

	cfg := config{
		steamAPIKey: steamAPIKey,
		steamID:     steamID,
	}

	if telegramNotifier {
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
		chatID := int64(telegramChatID)
		cfg.telegramBotToken = &telegramBotToken
		cfg.telegramChatID = &chatID
	}

	return cfg, nil
}
