package api

import (
	"encoding/json"
	"github.com/noah-blockchain/coin-price-backend/internal/usecase"
	"net/http"
)

// CoinHandler  represent the http handler for coin-history
type CoinHandler struct {
	app usecase.Usecase
}

// NewCoinHandler will initialize the / resources endpoint
func NewCoinPriceHandler(app usecase.Usecase) CoinHandler {
	handler := CoinHandler{
		app: app,
	}
	return handler
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
