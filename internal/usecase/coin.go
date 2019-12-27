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
	GetPrice(ctx context.Context, symbol string, date string, period string) ([]CoinPrice, error)
}

// Repository represent the coin's repository contract
type Repository interface {
	Store(ctx context.Context, coin *models.Coin) error
	GetBySymbol(ctx context.Context, symbol string) (*[]models.Coin, error)
	GetByDate(ctx context.Context, symbol string, start time.Time, end time.Time) (*[]models.Coin, error)
	GetSymbolNames(ctx context.Context) ([]string, error)
	GetLastPriceOnDate(ctx context.Context, symbol string, date time.Time) (*models.Coin, error)
	GetLastPrice(ctx context.Context, symbol string) (*models.Coin, error)
	GetLastPriceBeforeDate(ctx context.Context, symbol string, date time.Time) (*string, error)
}

// NewCoinUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewCoinUsecase(repo Repository) Usecase {
	return &app{
		repo: repo,
	}
}

type CoinPrice struct {
	Date  string `json:"date"`
	Price string `json:"value"`
}

func (a *app) GetPrice(c context.Context, symbol string, date string, period string) ([]CoinPrice, error) {
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
		days := end.Sub(start).Hours() / 24
		temp := make(map[string]string)
		keys := make([]string, int(days))
		for i := 0; i < int(days); i++ {
			start = end.AddDate(0, 0, -i)
			key := start.Format("02-01-2006")
			temp[key] = "0"
			keys[i] = key
		}

		coins, err := a.repo.GetByDate(c, symbol, start, end)
		if err != nil {
			return nil, err
		}
		for _, c := range *coins {
			key := c.CreatedAt.Format("02-01-2006")
			temp[key] = c.Price
		}
		res := make([]CoinPrice, int(days))
		for i, k := range keys {
			res[i].Date = k
			if temp[k] != "0" {
				res[i].Price = temp[k]
			} else {
				res[i].Price = temp[k]
				date, _ := time.Parse("02-01-2006", k)
				date.Add(23*60*60 + 59*60 + 59)
				p, _ := a.repo.GetLastPriceBeforeDate(c, symbol, date)
				if p != nil {
					res[i].Price = *p
				}
			}
		}

		return res, nil
	}
	coins, err := a.repo.GetBySymbol(c, symbol)
	if err != nil {
		return nil, err
	}

	res := make([]CoinPrice, len(*coins))
	for i, c := range *coins {
		res[i].Date = c.CreatedAt.Format("02-01-2006")
		res[i].Price = c.Price
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
