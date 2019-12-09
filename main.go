package main

import (
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	_coinHttpDeliver "github.com/noah-blockchain/coin-price-backend/coin/delivery/http"
	_coinRep "github.com/noah-blockchain/coin-price-backend/coin/repository"
	_coinUcase "github.com/noah-blockchain/coin-price-backend/coin/usecase"
	"github.com/noah-blockchain/coin-price-backend/middleware"
	"github.com/spf13/viper"
	"log"
	"time"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	dbConn, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbConn.Close()
	fmt.Println("Successfully connected!")

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	coinRepo := _coinRep.NewPsqlCoinRepository(dbConn.DB)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _coinUcase.NewCoinUsecase(coinRepo, timeoutContext)
	_coinHttpDeliver.NewCoinPriceHandler(e, au)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
