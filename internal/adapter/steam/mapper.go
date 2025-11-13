package steam

import (
	"fmt"

	"github.com/gobtronic/steam-purchase-notifier/internal/domain"
)

func gameDTOToDomain(dto gameDTO) domain.Game {
	return domain.Game{
		AppID:    dto.AppID,
		Name:     dto.Name,
		StoreURL: fmt.Sprintf("https://store.steampowered.com/app/%d", dto.AppID),
	}
}
