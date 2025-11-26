package gamestore

type storeDTO struct {
	Libraries []libraryDTO `json:"libraries"`
}

type libraryDTO struct {
	SteamID string `json:"steam_id"`
	AppIDs  []int  `json:"app_ids"`
}
