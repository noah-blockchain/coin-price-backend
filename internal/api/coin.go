package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/noah-blockchain/coin-price-backend/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CoinPrice struct {
	Date  string `json:"date"`
	Price string `json:"value"`
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

	respondWithJSON(w, getStatusCode(err), CoinPrice{cn.CreatedAt.Format("02-01-2006"), cn.Price})
}

// GetAllSymbolRecords will get all records for given symbol
func (a *CoinHandler) GetAllRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	period := r.URL.Query().Get("period")
	date := r.URL.Query().Get("date")

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	coinList, err := a.app.GetBySymbol(ctx, symbol, date, period)
	result := make([]CoinPrice, len(coinList))

	if err != nil {
		respondWithError(w, getStatusCode(err), err.Error())
		return
	}
	for i, c := range coinList {
		result[i].Price = c.Price
		result[i].Date = c.CreatedAt.Format("02-01-2006")
	}
	respondWithJSON(w, getStatusCode(err), result)
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
