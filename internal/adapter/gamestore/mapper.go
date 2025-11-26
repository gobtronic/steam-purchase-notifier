package gamestore

import "github.com/gobtronic/steam-purchase-notifier/internal/domain"

func libraryToDTO(library domain.Library) libraryDTO {
	var appIDs []int
	for _, g := range library.Games {
		appIDs = append(appIDs, g.AppID)
	}
	return libraryDTO{
		SteamID: library.SteamID,
		AppIDs:  appIDs,
	}
}

func storeDTOToDomain(dto storeDTO) []domain.Library {
	var users []domain.Library
	for _, lb := range dto.Libraries {
		var games []domain.Game
		for _, v := range lb.AppIDs {
			games = append(games, domain.Game{
				AppID: v,
			})
		}
		users = append(users, domain.Library{
			SteamID: lb.SteamID,
			Games:   games,
		})
	}
	return users
}
