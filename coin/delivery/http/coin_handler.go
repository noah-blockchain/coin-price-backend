package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/noah-blockchain/coin-price-backend/coin"
	"github.com/noah-blockchain/coin-price-backend/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type CoinHandler struct {
	CUsecase coin.Usecase
}

// NewCoinHandler will initialize the / resources endpoint
func NewCoinPriceHandler(e *echo.Echo, us coin.Usecase) {
	handler := &CoinHandler{
		CUsecase: us,
	}
	e.GET("/price", handler.GetCoinPrice)
}

type CoinPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// GetCoinPrice will get latest price for given symbol
func (a *CoinHandler) GetCoinPrice(c echo.Context) error {
	symbolP := c.QueryParam("symbol")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cn, err := a.CUsecase.GetLatestPrice(ctx, symbolP)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &CoinPrice{cn.Symbol, cn.Price})
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
