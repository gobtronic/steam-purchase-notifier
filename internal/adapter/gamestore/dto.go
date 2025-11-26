package gamestore

type storeDTO struct {
	Users []userDTO `json:"users"`
}

type userDTO struct {
	SteamID string `json:"steam_id"`
	AppIDs  []int  `json:"app_ids"`
}
