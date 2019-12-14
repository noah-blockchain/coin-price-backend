package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/noah-blockchain/coin-price-backend/internal/models"
	"github.com/noah-blockchain/coin-price-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type repo struct {
	Conn *sql.DB
}

// NewPsqlCoinRepository will create an object that represent the article.Repository interface
func NewPsqlCoinRepository(Conn *sql.DB) usecase.Repository {
	return &repo{Conn}
}

func (m *repo) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Coin, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.Coin, 0)
	for rows.Next() {
		t := new(models.Coin)
		err = rows.Scan(
			&t.ID,
			&t.Volume,
			&t.ReserveBalance,
			&t.Price,
			&t.Capitalization,
			&t.Symbol,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *repo) fetchByDate(ctx context.Context, query string, args ...interface{}) ([]*models.Coin, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.Coin, 0)
	for rows.Next() {
		t := new(models.Coin)
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

	return result, nil
}

func (m *repo) GetByID(ctx context.Context, id int64) (res *models.Coin, err error) {
	query := `SELECT id, volume, reserve_balance, price, capitalization, symbol, created_at FROM coins WHERE ID = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *repo) GetBySymbol(ctx context.Context, symbol string) (res []*models.Coin, err error) {
	query := `SELECT * FROM coins WHERE symbol = $1`

	list, err := m.fetch(ctx, query, symbol)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *repo) GetByDate(ctx context.Context, symbol string, start time.Time, end time.Time) (res []*models.Coin, err error) {
	query := `SELECT date_trunc('day', created_at) AS "day", AVG(price) 
				FROM coins WHERE symbol=$1 AND created_at>=$2 AND created_at<$3
				GROUP BY 1 
				ORDER BY 1`
	logrus.Info(start)
	logrus.Info(end)
	list, err := m.fetchByDate(ctx, query, symbol, start, end)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *repo) GetLatestPrice(ctx context.Context, symbol string) (res *models.Coin, err error) {
	query := `SELECT * FROM coins WHERE symbol = $1 ORDER BY created_at DESC LIMIT 1`

	list, err := m.fetch(ctx, query, symbol)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *repo) Store(ctx context.Context, c *models.Coin) error {
	query := `INSERT INTO public.coins(volume, reserve_balance, price, capitalization, symbol, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, c.Volume, c.ReserveBalance, c.Price, c.Capitalization, c.Symbol, c.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
