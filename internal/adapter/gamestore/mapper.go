package gamestore

import "github.com/gobtronic/steam-purchase-notifier/internal/domain"

func gamesToStoreDTO(games []domain.Game) storeDTO {
	var appIDs []int
	for _, g := range games {
		appIDs = append(appIDs, g.AppID)
	}
	return storeDTO{
		AppIDs: appIDs,
	}
}
