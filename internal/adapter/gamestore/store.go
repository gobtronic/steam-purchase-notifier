package gamestore

import (
	"encoding/json"
	"os"
	"path"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
)

type GameStore struct {
	filePath string
}

func NewGameStore() (*GameStore, error) {
	filePath := path.Join(os.Getenv("GOPATH"), "gamelist.json")
	return &GameStore{
		filePath: filePath,
	}, nil
}

func (s *GameStore) Write(games []domain.UserGames) error {
	dto, _ := s.read()
	for _, g := range games {
		cachedUserIndex := -1
		for i, u := range dto.Users {
			if u.SteamID == g.SteamID {
				cachedUserIndex = i
				break
			}
		}
		if cachedUserIndex >= 0 {
			dto.Users[cachedUserIndex] = gamesToUserDTO(g)
		} else {
			dto.Users = append(dto.Users, gamesToUserDTO(g))
		}
	}

	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(dto); err != nil {
		return err
	}
	return nil
}

func (s *GameStore) Read() ([]domain.UserAppIDs, error) {
	dto, err := s.read()
	if err != nil {
		return []domain.UserAppIDs{}, err
	}
	return userDTOsToDomain(dto), nil
}

func (s *GameStore) read() (storeDTO, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return storeDTO{}, err
	}
	defer file.Close()

	var dto storeDTO
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dto); err != nil {
		return storeDTO{}, err
	}
	return dto, nil
}
