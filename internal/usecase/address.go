package usecase

import (
	"context"
	"time"

	"github.com/noah-blockchain/coin-price-backend/internal/models"
	coin_extender "github.com/noah-blockchain/coinExplorer-tools"
)

func (a *app) CreateAddressHistory(ctx context.Context, address coin_extender.Address) error {
	return a.repo.StoreAddress(ctx, &models.Address{
		Address:   address.Address,
		Symbol:    address.Symbol,
		Amount:    address.Amount,
		CreatedAt: time.Unix(address.CreatedAt.Seconds, int64(address.CreatedAt.Nanos)),
	})
}
