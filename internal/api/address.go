package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noah-blockchain/coin-price-backend/internal/models"
	"github.com/noah-blockchain/coinExplorer-tools/helpers"
)

// Get price of coin for date range
func (a *CoinHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	if len(address) != 45 {
		respondWithError(w, http.StatusBadRequest, models.ErrWrongNoahAddress.Error())
		return
	}

	period := r.URL.Query().Get("period")
	date := r.URL.Query().Get("date")
	address = helpers.RemovePrefixFromAddress(address)
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	balances, err := a.app.GetAddressBalance(ctx, address, date, period)
	if err != nil {
		respondWithError(w, getStatusCode(err), err.Error())
		return
	}

	respondWithJSON(w, getStatusCode(err), balances)
}
