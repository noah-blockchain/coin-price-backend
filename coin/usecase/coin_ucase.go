package usecase

import (
	"context"
	"github.com/noah-blockchain/coin-price-backend/coin"
	"github.com/noah-blockchain/coin-price-backend/models"
	"time"
)

type coinUsecase struct {
	coinRepo       coin.Repository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewCoinUsecase(a coin.Repository, timeout time.Duration) coin.Usecase {
	return &coinUsecase{
		coinRepo:       a,
		contextTimeout: timeout,
	}
}

func (a *coinUsecase) GetLatestPrice(c context.Context, symbol string) (*models.Coin, error) {

	res, err := a.coinRepo.GetLatestPrice(c, symbol)
	if err != nil {
		return nil, err
	}

	return res, nil
}
