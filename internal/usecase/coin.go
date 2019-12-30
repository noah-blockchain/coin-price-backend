package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/noah-blockchain/coin-price-backend/internal/helpers"
	"github.com/noah-blockchain/coin-price-backend/internal/models"
	coin_extender "github.com/noah-blockchain/coinExplorer-tools"
)

const (
	precision        = 100
	endOfDayDuration = (23*60*60 + 59*60 + 59) * time.Second
	timeParseLayout  = "02-01-2006"
)

type app struct {
	repo Repository
}

// Usecase represent the coin's usecases
type Usecase interface {
	CreateCoinInfo(ctx context.Context, coin coin_extender.Coin) error
	CreateAddressHistory(ctx context.Context, address coin_extender.Address) error
	GetPrice(ctx context.Context, symbol string, date string, period string) ([]CoinPrice, error)
	GetAddressBalance(ctx context.Context, address string, date string, period string) ([]Balance, error)
}

// Repository represent the coin's repository contract
type Repository interface {
	Store(ctx context.Context, coin *models.Coin) error
	StoreAddress(ctx context.Context, address *models.Address) error
	GetBySymbol(ctx context.Context, symbol string) (*[]models.Coin, error)
	GetByDate(ctx context.Context, symbol string, start time.Time, end time.Time) (*[]models.Coin, error)
	GetSymbolNames(ctx context.Context) ([]string, error)
	GetLastPriceOnDate(ctx context.Context, symbol string, date time.Time) (*models.Coin, error)
	GetLastPrice(ctx context.Context, symbol string) (*models.Coin, error)
	GetLastPriceBeforeDate(ctx context.Context, symbol string, date time.Time) (string, error)
	GetAddressBalances(ctx context.Context, address string, date time.Time) ([]models.Address, error)
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

type Balance struct {
	Date    string `json:"date"`
	Balance string `json:"balance"`
}

func (a *app) GetPrice(c context.Context, symbol string, date string, period string) ([]CoinPrice, error) {
	if date != "" || period != "" {
		end, err := a.parseAndGetEndOfDay(date)
		if err != nil {
			return nil, err
		}

		start, err := a.getStartOfPeriod(end, period)
		if err != nil {
			return nil, err
		}

		days := end.Sub(start).Hours() / 24
		temp := make(map[string]string)
		keys := make([]string, int(days))
		for i := 0; i < int(days); i++ {
			start = end.AddDate(0, 0, -i)
			key := start.Format(timeParseLayout)
			temp[key] = "0"
			keys[i] = key
		}

		coins, err := a.repo.GetByDate(c, symbol, start, end)
		if err != nil {
			return nil, err
		}
		for _, c := range *coins {
			key := c.CreatedAt.Format(timeParseLayout)
			temp[key] = c.Price
		}
		res := make([]CoinPrice, int(days))
		for i, k := range keys {
			res[i].Date = k
			res[i].Price = temp[k]
			if temp[k] == "0" {
				date, _ := a.parseAndGetEndOfDay(k)
				p, err := a.repo.GetLastPriceBeforeDate(c, symbol, date)
				if err == nil {
					res[i].Price = p
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
		res[i].Date = c.CreatedAt.Format(timeParseLayout)
		res[i].Price = c.Price
	}

	return res, nil
}

func (a *app) GetAddressBalance(c context.Context, address string, date string, period string) ([]Balance, error) {
	if date != "" || period != "" {
		end, err := a.parseAndGetEndOfDay(date)
		if err != nil {
			return nil, err
		}

		start, err := a.getStartOfPeriod(end, period)
		if err != nil {
			return nil, err
		}

		days := end.Sub(start).Hours() / 24
		temp := make(map[string]string)
		keys := make([]string, int(days))
		for i := 0; i < int(days); i++ {
			start = end.AddDate(0, 0, -i)
			key := start.Format(timeParseLayout)
			temp[key] = "0"
			keys[i] = key
		}
		for _, key := range keys {
			d, _ := a.parseAndGetEndOfDay(key)
			balances, err := a.repo.GetAddressBalances(c, address, d)
			if err != nil {
				return nil, err
			}
			if balances != nil {
				sum := helpers.NewFloat(0, precision)
				for _, b := range balances {
					price, err := a.repo.GetLastPriceBeforeDate(c, b.Symbol, d)
					if err == nil {
						priceFloat, _ := helpers.NewFloat(0, precision).SetString(price)
						amount, _ := helpers.NewFloat(0, precision).SetString(b.Amount)
						sum.Add(sum, priceFloat.Mul(priceFloat, amount))
					}
				}
				temp[key] = sum.String()
			}
		}
		balances := make([]Balance, len(temp))
		for i, k := range keys {
			balances[i].Date = k
			balances[i].Balance = temp[k]
		}
		return balances, nil
	}
	return nil, models.ErrBadParamInput
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

func (a *app) parseAndGetEndOfDay(date string) (time.Time, error) {
	end, err := time.Parse(timeParseLayout, date)
	if err != nil {
		return time.Now(), errors.New(fmt.Sprintf("Failed to parse date format : %s", date))
	}
	return end.Add(endOfDayDuration), nil // date must be with 23:59:59 time in the end
}

func (a *app) getStartOfPeriod(end time.Time, period string) (time.Time, error) {
	var start time.Time
	switch period {
	case "WEEK":
		start = end.AddDate(0, 0, -7)
	case "MONTH":
		start = end.AddDate(0, -1, 0)
	case "YEAR":
		start = end.AddDate(-1, 0, 0)
	default:
		return time.Now(), errors.New(fmt.Sprintf("Incorrect format : %s", period))
	}
	return start, nil
}
