package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

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
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		var notifiers []port.Notifier
		if telegramNotifier {
			telegramCfg, err := telegram.LoadConfig()
			if err != nil {
				return err
			}
			notifier := telegram.NewTelegramNotifier(telegramCfg)
			notifiers = append(notifiers, notifier)
		}

		steamClient := steam.NewSteamClient(cfg.steamAPIKey, cfg.steamID, http.DefaultClient)
		gameStoreFilePath := path.Join(os.Getenv("GOPATH"), "gamelist.json")
		gameStore := gamestore.NewGameStore(gameStoreFilePath)
		games, err := steamClient.FetchGames()
		if err != nil {
			return err
		}
		newGames, _ := usecase.FilterNewGames(games, gameStore)
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
	steamAPIKey string
	steamID     string
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

	return cfg, nil
}
