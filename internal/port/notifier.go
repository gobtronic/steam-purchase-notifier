package port

import "github.com/gobtronic/steam-purchase-notifier/internal/domain"

type Notifier interface {
	Notify(game domain.Game) error
}
