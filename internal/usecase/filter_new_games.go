package usecase

import (
	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
	"github.com/gobtronic/steam-purchase-notifier/internal/port"
)

func FilterNewGames(games []domain.Game, gameStore port.GameStore) ([]domain.Game, error) {
	cachedIDs, err := gameStore.Read()
	if err != nil {
		return []domain.Game{}, err
	}

	var freshIDs []int
	for _, g := range games {
		freshIDs = append(freshIDs, g.AppID)
	}

	flaggedIDs := make(map[int]bool)
	for _, f := range freshIDs {
		flaggedIDs[f] = true
	}

	for _, c := range cachedIDs {
		if _, ok := flaggedIDs[c]; ok {
			flaggedIDs[c] = false
		}
	}

	var newGames []domain.Game
	for _, g := range games {
		if flaggedIDs[g.AppID] {
			newGames = append(newGames, g)
		}
	}

	return newGames, nil
}
