package coin

import (
	"context"

	"github.com/noah-blockchain/coin-price-backend/models"
)

// Usecase represent the coin's usecases
type Usecase interface {
	GetLatestPrice(ctx context.Context, symbol string) (*models.Coin, error)
}
