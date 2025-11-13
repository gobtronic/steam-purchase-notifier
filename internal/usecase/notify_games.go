package usecase

import (
	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
	"github.com/gobtronic/steam-purchase-notifier/internal/port"
)

func NotifyGames(games []domain.Game, notifier port.Notifier) {
	for _, g := range games {
		notifier.Notify(g)
	}
}
