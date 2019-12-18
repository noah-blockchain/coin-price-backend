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

func (m *repo) GetByDate(ctx context.Context, symbol string, start time.Time, end time.Time) (*[]models.Coin, error) {
	query := `SELECT date_trunc('day', created_at) AS "day", AVG(price) 
				FROM coins WHERE symbol=$1 AND created_at>=$2 AND created_at<$3
				GROUP BY 1 
				ORDER BY 1`
	list, err := m.fetchByDate(ctx, query, symbol, start, end)
	if err != nil || list == nil {
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
