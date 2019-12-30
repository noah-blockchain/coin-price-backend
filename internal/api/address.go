package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// Get price of coin for date range
func (a *CoinHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	period := r.URL.Query().Get("period")
	date := r.URL.Query().Get("date")
	address = removeNoahPrefix(address)
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

func removeNoahPrefix(raw string) string {
	return raw[5:]
}
