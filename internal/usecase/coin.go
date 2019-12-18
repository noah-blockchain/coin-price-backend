package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/noah-blockchain/coin-price-backend/internal/models"
	coin_extender "github.com/noah-blockchain/coinExplorer-tools"
)

type app struct {
	repo Repository
}

// Usecase represent the coin's usecases
type Usecase interface {
	CreateCoinInfo(ctx context.Context, coin coin_extender.Coin) error
	GetBySymbol(ctx context.Context, symbol string, date string, period string) (*[]models.Coin, error)
}

// Repository represent the coin's repository contract
type Repository interface {
	Store(ctx context.Context, coin *models.Coin) error
	GetBySymbol(ctx context.Context, symbol string) (*[]models.Coin, error)
	GetByDate(ctx context.Context, symbol string, start time.Time, end time.Time) (*[]models.Coin, error)
}

// NewCoinUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewCoinUsecase(repo Repository) Usecase {
	return &app{
		repo: repo,
	}
}

func (a *app) GetBySymbol(c context.Context, symbol string, date string, period string) (*[]models.Coin, error) {
	if date != "" || period != "" {
		layout := "02-01-2006"
		end, err := time.Parse(layout, date)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to parse date format : %s", date))
		}
		end = end.Add((23*60*60 + 59*60 + 59) * time.Second) // date must be with 23:59:59 time in the end

		var start time.Time
		switch period {
		case "WEEK":
			start = end.AddDate(0, 0, -7)
		case "MONTH":
			start = end.AddDate(0, -1, 0)
		case "YEAR":
			start = end.AddDate(-1, 0, 0)
		default:
			return nil, errors.New(fmt.Sprintf("Incorrect format : %s", period))
		}
		res, err := a.repo.GetByDate(c, symbol, start, end)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	res, err := a.repo.GetBySymbol(c, symbol)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *app) CreateCoinInfo(ctx context.Context, coin coin_extender.Coin) error {
	return a.repo.Store(ctx, &models.Coin{
		Symbol:         coin.Symbol,
		Price:          coin.Price,
		Capitalization: coin.Capitalization,
		ReserveBalance: coin.ReserveBalance,
		Volume:         coin.Volume,
		CreatedAt:      time.Unix(coin.CreatedAt.Seconds, int64(coin.CreatedAt.Nanos)),
	})
}
