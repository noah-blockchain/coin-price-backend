package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/noah-blockchain/coin-price-backend/internal/models"
	"github.com/noah-blockchain/coin-price-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
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

func (m *repo) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Coin, string, error) {
	query := `SELECT id, volume, reserve_balance, price, capitalization, symbol, created_at FROM coins WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}

	res, err := m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return res, nextCursor, err
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
				FROM coins WHERE symbol=$1 AND created_at>=$2 AND created_at<=$3
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
	query := `INSERT INTO public.coins(volume, reserve_balance, price, capitalization, symbol)
	VALUES ($1, $2, $3, $4, $5)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, c.Volume, c.ReserveBalance, c.Price, c.Capitalization, c.Symbol)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = uint64(lastID)
	return nil
}

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
