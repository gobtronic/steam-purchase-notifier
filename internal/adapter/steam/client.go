package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
)

const STEAM_API_BASE_URL = "https://api.steampowered.com/"

type SteamClient struct {
	apiKey   string
	SteamIDs []string
	client   *http.Client
}

func NewSteamClient(client *http.Client) (*SteamClient, error) {
	apiKey := os.Getenv("STEAM_API_KEY")
	if apiKey == "" {
		return &SteamClient{}, fmt.Errorf("Please set the STEAM_API_KEY environment variable with your Steam API key")
	}
	steamID := os.Getenv("STEAM_IDS")
	if steamID == "" {
		return &SteamClient{}, fmt.Errorf("Please set the STEAM_IDS environment variable with Steam IDs that will be watched (separated by commas)")
	}
	steamIDs := strings.Split(steamID, ",")
	return &SteamClient{
		apiKey:   apiKey,
		SteamIDs: steamIDs,
		client:   client,
	}, nil
}

func (c *SteamClient) FetchGames(userID string) ([]domain.Game, error) {
	u, err := url.Parse(STEAM_API_BASE_URL)
	if err != nil {
		return []domain.Game{}, err
	}
	u.Path = path.Join(u.Path, "IPlayerService/GetOwnedGames/v1")
	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("steamid", userID)
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
