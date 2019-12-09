package coin

import (
	"context"

	"github.com/noah-blockchain/coin-price-backend/models"
)

// Repository represent the coin's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Coin, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.Coin, error)
	GetLatestPrice(ctx context.Context, symbol string) (*models.Coin, error)
}
