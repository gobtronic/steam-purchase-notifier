package domain

type Game struct {
	AppID    int
	Name     string
	StoreURL string
}

type Library struct {
	SteamID string
	Games   []Game
}
