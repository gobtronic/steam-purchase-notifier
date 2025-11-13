package steam

type ownedGamesResponseDTO struct {
	ResponseDTO struct {
		GameCount int       `json:"game_count"`
		GamesDTOs []gameDTO `json:"games"`
	} `json:"response"`
}

type gameDTO struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
