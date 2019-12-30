package repository

import (
	"context"
	"time"

	"github.com/noah-blockchain/coin-price-backend/internal/models"
)

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
