package usecase

import (
	"fmt"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
	"github.com/gobtronic/steam-purchase-notifier/internal/port"
)

func FilterNewGames(library domain.Library, gameStore port.GameStore) ([]domain.Game, error) {
	cache, err := gameStore.Read()
	if err != nil {
		return []domain.Game{}, err
	}

	var cachedLibrary *domain.Library
	for _, lb := range cache {
		if lb.SteamID == library.SteamID {
			cachedLibrary = &lb
		}
	}
	if cachedLibrary == nil {
		return []domain.Game{}, fmt.Errorf("Could not find cached library for Steam ID %s", library.SteamID)
	}

	var freshIDs []int
	for _, g := range library.Games {
		freshIDs = append(freshIDs, g.AppID)
	}

	flaggedIDs := make(map[int]bool)
	for _, f := range freshIDs {
		flaggedIDs[f] = true
	}

	for _, g := range cachedLibrary.Games {
		if _, ok := flaggedIDs[g.AppID]; ok {
			flaggedIDs[g.AppID] = false
		}
	}

	var newGames []domain.Game
	for _, g := range library.Games {
		if flaggedIDs[g.AppID] {
			newGames = append(newGames, g)
		}
	}

	return newGames, nil
}
