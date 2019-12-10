package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noah-blockchain/coin-price-backend/internal/models"
	"github.com/sirupsen/logrus"
)

type CoinPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// GetCoinPrice will get latest price for given symbol
func (a *CoinHandler) GetCoinPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cn, err := a.app.GetLatestPrice(ctx, symbol)
	if err != nil {
		respondWithError(w, getStatusCode(err), err.Error())
		return
	}

	respondWithJSON(w, getStatusCode(err), CoinPrice{cn.Symbol, cn.Price})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
