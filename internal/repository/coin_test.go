package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/noah-blockchain/coin-price-backend/internal/config"
	"github.com/noah-blockchain/coin-price-backend/internal/env"
	"github.com/noah-blockchain/coin-price-backend/internal/usecase"
	"time"

	"log"
	"os"
	"testing"
)

var rep usecase.Repository

func TestMain(m *testing.M) {
	var cfg = config.Config{}
	cfg.DbPort = env.GetEnvAsInt("DB_PORT", 5432)
	cfg.DbHost = os.Getenv("DB_HOST")
	cfg.DbUser = os.Getenv("DB_USER")
	cfg.DbName = os.Getenv("DB_NAME")
	cfg.DbPass = os.Getenv("DB_PASSWORD")

	dbDsnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName,
	)
	var err error
	dbConn, err := sqlx.Connect("postgres", dbDsnString)
	if err != nil {
		log.Panicln(err)
	}
	defer dbConn.Close()
	fmt.Println("DB connected successful!")
	rep = NewPsqlCoinRepository(dbConn.DB)

	code := m.Run()

	os.Exit(code)
}

func TestGetByDate(t *testing.T) {
	date := "11-12-2019"
	end, err := time.Parse("02-01-2006", date)
	if err != nil {
		t.Error(err)
	}
	start := end.AddDate(0, 0, -1)
	// should return avg coin price for the day 10-12-2019
	res, err := rep.GetByDate(context.TODO(), "NOAH", start, end)
	if err != nil {
		t.Error(err)
	}
	if len(res) > 0 {
		if len(res) != 1 {
			t.Errorf("Should return result for one day")
		}
		if res[0].CreatedAt.Format("02-01-2006") != "10-12-2019" {
			t.Errorf("Should return result for day 10-12-2019 ")
		}
	}
}
