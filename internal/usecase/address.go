package usecase

import (
	"context"
	"time"

	"github.com/noah-blockchain/coin-price-backend/internal/helpers"
	"github.com/noah-blockchain/coin-price-backend/internal/models"
	coin_extender "github.com/noah-blockchain/coinExplorer-tools"
)

type Balance struct {
	Date    string `json:"date"`
	Balance string `json:"balance"`
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

func (a *app) CreateAddressHistory(ctx context.Context, address coin_extender.Address) error {
	return a.repo.StoreAddress(ctx, &models.Address{
		Address:   address.Address,
		Symbol:    address.Symbol,
		Amount:    address.Amount,
		CreatedAt: time.Unix(address.CreatedAt.Seconds, int64(address.CreatedAt.Nanos)),
	})
}
