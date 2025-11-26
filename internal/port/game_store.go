package port

import "github.com/gobtronic/steam-purchase-notifier/internal/domain"

type GameStore interface {
	Write(games []domain.Library) error
	Read() ([]domain.Library, error)
}
