package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/noah-blockchain/coin-price-backend/internal/models"
	"github.com/noah-blockchain/coin-price-backend/internal/usecase"
)

type repo struct {
	db *sqlx.DB
}

// NewPsqlCoinRepository will create an object that represent the article.Repository interface
func NewPsqlCoinRepository(db *sqlx.DB) usecase.Repository {
	return &repo{db}
}

func (m *repo) GetBySymbol(ctx context.Context, symbol string) (*[]models.Coin, error) {
	var history []models.Coin
	var err error
	err = m.db.Select(&history, "SELECT * FROM public.coins WHERE symbol = $1", symbol)
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (m *repo) GetSymbolNames(ctx context.Context) ([]string, error) {
	var symbolNames []string
	var err error
	err = m.db.Select(&symbolNames, "SELECT DISTINCT(symbol) FROM coins")
	if err != nil {
		return nil, err
	}
	return symbolNames, nil
}

func (m *repo) GetLastPriceBeforeDate(ctx context.Context, symbol string, date time.Time) (*string, error) {
	var lastPrice string
	var err error

	stmt, err := m.db.Preparex(`SELECT last(price, created_at)
									  FROM coins
									  WHERE symbol = $1 AND created_at < $2`)

	if err != nil {
		return nil, err
	}
	err = stmt.Get(&lastPrice, symbol, date)
	if err != nil {
		return nil, err
	}
	return &lastPrice, nil
}

func (m *repo) GetAddressBalances(ctx context.Context, address string, date time.Time) ([]models.Address, error) {
	var balances []models.Address
	var err error

	stmt, err := m.db.Preparex(`SELECT *  FROM addresses 
						WHERE address=$1
						AND date_trunc('day', created_at) = 
						(SELECT date_trunc('day', created_at) as lastday 
						 FROM addresses WHERE created_at < $2
						 LIMIT 1
						)
							`)

	if err != nil {
		return nil, err
	}
	err = stmt.Select(&balances, address, date)
	if err != nil {
		return nil, err
	}
	return balances, nil
}

func (m *repo) GetLastPriceOnDate(ctx context.Context, symbol string, date time.Time) (*models.Coin, error) {
	var coin models.Coin
	var err error
	stmt, err := m.db.Preparex(`SELECT * from coins 
									WHERE symbol=$1 AND date_trunc('day', created_at) = $2 
									ORDER BY created_at DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	err = stmt.Get(&coin, symbol, date)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}

func (m *repo) GetLastPrice(ctx context.Context, symbol string) (*models.Coin, error) {
	var coin models.Coin
	var err error
	stmt, err := m.db.Preparex(`SELECT * FROM coins 
 									  WHERE symbol = $1
									  ORDER BY created_at DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	err = stmt.Get(&coin, symbol)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}

func (m *repo) GetByDate(ctx context.Context, symbol string, start time.Time, end time.Time) (*[]models.Coin, error) {
	query := `SELECT date_trunc('day', created_at) AS "day", AVG(price) 
				FROM coins WHERE symbol=$1 AND created_at>=$2 AND created_at<$3
				GROUP BY 1 
				ORDER BY 1`
	list, err := m.fetchByDate(ctx, query, symbol, start, end)
	if err != nil {
		return nil, err
	}

	if len(*list) > 0 {
		return list, nil
	}
	return nil, models.ErrNotFound
}

func (m *repo) Store(ctx context.Context, c *models.Coin) error {
	query := `INSERT INTO public.coins(volume, reserve_balance, price, capitalization, symbol, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, c.Volume, c.ReserveBalance, c.Price, c.Capitalization, c.Symbol, c.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *repo) StoreAddress(ctx context.Context, address *models.Address) error {
	query := `INSERT INTO public.addresses(address, symbol, amount, created_at)
	VALUES ($1, $2, $3, $4)`
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, address.Address, address.Symbol, address.Amount, address.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *repo) fetchByDate(ctx context.Context, query string, args ...interface{}) (*[]models.Coin, error) {
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]models.Coin, 0)
	for rows.Next() {
		t := models.Coin{}
		err = rows.Scan(
			&t.CreatedAt,
			&t.Price,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return &result, nil
}
