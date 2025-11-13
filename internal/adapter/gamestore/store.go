package gamestore

import (
	"encoding/json"
	"os"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
)

type GameStore struct {
	filePath string
}

func NewGameStore(filePath string) *GameStore {
	return &GameStore{
		filePath: filePath,
	}
}

func (s *GameStore) Write(games []domain.Game) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	dto := gamesToStoreDTO(games)
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(dto); err != nil {
		return err
	}
	return nil
}

func (s *GameStore) Read() ([]int, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return []int{}, err
	}
	defer file.Close()

	var dto storeDTO
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dto); err != nil {
		return []int{}, err
	}
	return dto.AppIDs, err
}
