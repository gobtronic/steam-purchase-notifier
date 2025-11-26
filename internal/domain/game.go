package domain

type Game struct {
	AppID    int
	Name     string
	StoreURL string
}

type UserGames struct {
	SteamID string
	Games   []Game
}

type UserAppIDs struct {
	SteamID string
	AppIDs  []int
}
