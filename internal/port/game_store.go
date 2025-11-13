package port

import "github.com/gobtronic/steam-purchase-notifier/internal/domain"

type GameStore interface {
	Write(games []domain.Game) error
	Read() ([]int, error)
}
