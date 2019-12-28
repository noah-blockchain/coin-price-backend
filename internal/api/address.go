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

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	coinList, err := a.app.GetPrice(ctx, address, date, period)
	if err != nil {
		respondWithError(w, getStatusCode(err), err.Error())
		return
	}

	respondWithJSON(w, getStatusCode(err), coinList)
}
