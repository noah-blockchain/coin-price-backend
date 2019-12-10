package usecase

import (
	"context"

	"github.com/noah-blockchain/coin-price-backend/internal/models"
)

type app struct {
	repo Repository
}

// Usecase represent the coin's usecases
type Usecase interface {
	GetLatestPrice(ctx context.Context, symbol string) (*models.Coin, error)
}

// Repository represent the coin's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Coin, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.Coin, error)
	GetLatestPrice(ctx context.Context, symbol string) (*models.Coin, error)
	Store(ctx context.Context, coin *models.Coin) error
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewCoinUsecase(repo Repository) Usecase {
	return &app{
		repo: repo,
	}
}

func (a *app) GetLatestPrice(c context.Context, symbol string) (*models.Coin, error) {
	res, err := a.repo.GetLatestPrice(c, symbol)
	if err != nil {
		return nil, err
	}

	return res, nil
}
