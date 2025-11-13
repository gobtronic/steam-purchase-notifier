package port

import "github.com/gobtronic/steam-purchase-notifier/internal/domain"

type SteamClient interface {
	FetchGames() ([]domain.Game, error)
}
