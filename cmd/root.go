package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/discord"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/gamestore"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/steam"
	"github.com/gobtronic/steam-purchase-notifier/internal/adapter/telegram"
	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
	"github.com/gobtronic/steam-purchase-notifier/internal/port"
	"github.com/gobtronic/steam-purchase-notifier/internal/usecase"
	"github.com/spf13/cobra"
)

var telegramNotifier bool
var discordNotifier bool
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

		if discordNotifier {
			notifier, err := discord.NewDiscordNotifier()
			if err != nil {
				log.Fatal(err)
			}
			notifiers = append(notifiers, notifier)
		}

		gameStore, err := gamestore.NewGameStore()
		if err != nil {
			log.Fatal(err)
		}

		var libraries []domain.Library
		for _, steamID := range steamClient.SteamIDs {
			games, err := steamClient.FetchGames(steamID)
			if err != nil {
				log.Println(err)
				continue
			}
			library := domain.Library{
				SteamID: steamID,
				Games:   games,
			}
			libraries = append(libraries, library)
			newGames, _ := usecase.FilterNewGames(library, gameStore)
			if len(newGames) > 0 && len(games) != len(newGames) {
				for _, notifier := range notifiers {
					usecase.NotifyGames(newGames, notifier)
				}
			}

		}
		gameStore.Write(libraries)
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
	rootCmd.Flags().BoolVar(&discordNotifier, "discord", false, "discord")
}
