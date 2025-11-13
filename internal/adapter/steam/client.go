package steam

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
)

const STEAM_API_BASE_URL = "https://api.steampowered.com/"

type SteamClient struct {
	apiKey  string
	steamID string
	client  *http.Client
}

func NewSteamClient(apiKey, steamID string, client *http.Client) *SteamClient {
	return &SteamClient{
		apiKey:  apiKey,
		steamID: steamID,
		client:  client,
	}
}

func (c *SteamClient) FetchGames() ([]domain.Game, error) {
	u, err := url.Parse(STEAM_API_BASE_URL)
	if err != nil {
		return []domain.Game{}, err
	}
	u.Path = path.Join(u.Path, "IPlayerService/GetOwnedGames/v1")
	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("steamid", c.steamID)
	params.Set("include_appinfo", "true")
	params.Set("include_played_free_games", "true")
	u.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return []domain.Game{}, err
	}
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []domain.Game{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []domain.Game{}, err
	}
	var gamesDTO ownedGamesResponseDTO
	err = json.Unmarshal(body, &gamesDTO)
	if err != nil {
		return []domain.Game{}, err
	}

	var games []domain.Game
	for _, dto := range gamesDTO.ResponseDTO.GamesDTOs {
		games = append(games, gameDTOToDomain(dto))
	}

	return games, nil
}
