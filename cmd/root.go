package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/gamestore"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/steam"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/telegram"
	"github.com/gobtronic/steam-purchase-notifier/internal/port"
	"github.com/gobtronic/steam-purchase-notifier/internal/usecase"
	"github.com/spf13/cobra"
)

var telegramNotifier bool
var rootCmd = &cobra.Command{
	Use:   "steam-purchase-notifier",
	Short: "Watch a Steam account purchases through notifications",
	Run: func(cmd *cobra.Command, args []string) {
		steamClient, err := steam.NewSteamClient(http.DefaultClient)
		if err != nil {
			log.Fatal(err)
		}

		var notifiers []port.Notifier
		if telegramNotifier {
			notifier, err := telegram.NewTelegramNotifier()
			if err != nil {
				log.Fatal(err)
			}
			notifiers = append(notifiers, notifier)
		}

		gameStore, err := gamestore.NewGameStore()
		if err != nil {
			log.Fatal(err)
		}
		games, err := steamClient.FetchGames()
		if err != nil {
			log.Fatal(err)
		}
		newGames, _ := usecase.FilterNewGames(games, gameStore)
		gameStore.Write(games)

		if len(games) != len(newGames) {
			for _, notifier := range notifiers {
				usecase.NotifyGames(newGames, notifier)
			}
		}
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
