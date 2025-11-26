package gamestore

import "github.com/gobtronic/steam-purchase-notifier/internal/domain"

func gamesToUserDTO(ug domain.UserGames) userDTO {
	var appIDs []int
	for _, g := range ug.Games {
		appIDs = append(appIDs, g.AppID)
	}
	return userDTO{
		SteamID: ug.SteamID,
		AppIDs:  appIDs,
	}
}

func userDTOsToDomain(dto storeDTO) []domain.UserAppIDs {
	var users []domain.UserAppIDs
	for _, u := range dto.Users {
		users = append(users, domain.UserAppIDs{
			SteamID: u.SteamID,
			AppIDs:  u.AppIDs,
		})
	}
	return users
}
