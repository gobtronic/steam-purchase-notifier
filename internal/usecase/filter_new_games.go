package usecase

import (
	"fmt"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
	"github.com/gobtronic/steam-purchase-notifier/internal/port"
)

func FilterNewGames(userGames domain.UserGames, gameStore port.GameStore) ([]domain.Game, error) {
	cache, err := gameStore.Read()
	if err != nil {
		return []domain.Game{}, err
	}

	var cachedUserAppIDs *domain.UserAppIDs
	for _, u := range cache {
		if u.SteamID == userGames.SteamID {
			cachedUserAppIDs = &u
		}
	}
	if cachedUserAppIDs == nil {
		return []domain.Game{}, fmt.Errorf("Could not find cached app IDs for Steam ID %s", userGames.SteamID)
	}

	var freshIDs []int
	for _, g := range userGames.Games {
		freshIDs = append(freshIDs, g.AppID)
	}

	flaggedIDs := make(map[int]bool)
	for _, f := range freshIDs {
		flaggedIDs[f] = true
	}

	for _, c := range cachedUserAppIDs.AppIDs {
		if _, ok := flaggedIDs[c]; ok {
			flaggedIDs[c] = false
		}
	}

	var newGames []domain.Game
	for _, g := range userGames.Games {
		if flaggedIDs[g.AppID] {
			newGames = append(newGames, g)
		}
	}

	return newGames, nil
}
