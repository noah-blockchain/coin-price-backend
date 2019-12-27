package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/noah-blockchain/coin-price-backend/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Get price of coin for date range
func (a *CoinHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	period := r.URL.Query().Get("period")
	date := r.URL.Query().Get("date")

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	coinList, err := a.app.GetPrice(ctx, symbol, date, period)
	if err != nil {
		respondWithError(w, getStatusCode(err), err.Error())
		return
	}

	respondWithJSON(w, getStatusCode(err), coinList)
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
