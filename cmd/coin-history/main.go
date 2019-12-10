package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/noah-blockchain/coin-price-backend/internal/api"
	"github.com/noah-blockchain/coin-price-backend/internal/config"
	"github.com/noah-blockchain/coin-price-backend/internal/env"
	"github.com/noah-blockchain/coin-price-backend/internal/repository"
	"github.com/noah-blockchain/coin-price-backend/internal/usecase"
)

var cfg = config.Config{}

func init() {
	flag.IntVar(&cfg.DbPort, "db.port", env.GetEnvAsInt("DB_PORT", 5432), "db port not exist")
	flag.StringVar(&cfg.DbHost, "db.host", os.Getenv("DB_HOST"), "db host not exist")
	flag.StringVar(&cfg.DbUser, "db.user", os.Getenv("DB_USER"), "db user not exist")
	flag.StringVar(&cfg.DbName, "db.name", os.Getenv("DB_NAME"), "db name not exist")
	flag.StringVar(&cfg.DbPass, "db.pass", os.Getenv("DB_PASSWORD"), "db pass not exist")
	flag.IntVar(&cfg.ServicePort, "service_port", env.GetEnvAsInt("SERVICE_PORT", 10500), "service port not exist")
	flag.BoolVar(&cfg.Debug, "debug", env.GetEnvAsBool("DEBUG", true), "debug not exist")
}

func main() {
	flag.Parse()
	switch {
	case cfg.DbHost == "":
		log.Panicf("Invalid value %s for field %s", cfg.DbHost, "db.host")
	case cfg.DbUser == "":
		log.Panicf("Invalid value %s for field %s", cfg.DbUser, "db.user")
	case cfg.DbName == "":
		log.Panicf("Invalid value %s for field %s", cfg.DbName, "db.name")
	case cfg.DbPass == "":
		log.Panicf("Invalid value %s for field %s", cfg.DbPass, "db.pass")
	case cfg.ServicePort <= 0:
		log.Panicf("Invalid value %d for field %s", cfg.ServicePort, "service_port")
	case cfg.DbPort <= 0:
		log.Panicf("Invalid value %d for field %s", cfg.DbPort, "db.port")

	}

	dbDsnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName,
	)
	dbConn, err := sqlx.Connect("postgres", dbDsnString)
	if err != nil {
		log.Panicln(err)
	}
	defer dbConn.Close()
	fmt.Println("DB connected successful!")

	driver, err := postgres.WithInstance(dbConn.DB, &postgres.Config{})
	err = runMigrations(driver)
	if err != nil {
		log.Panicln(err)
	}

	repo := repository.NewPsqlCoinRepository(dbConn.DB)
	app := usecase.NewCoinUsecase(repo)
	handler := api.NewCoinPriceHandler(app)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/price/{symbol}", handler.GetCoinPrice).Methods("GET")

	fmt.Println("Starting coin-history service with port", cfg.ServicePort)
	log.Panicln(http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServicePort), router))
}

func runMigrations(driver database.Driver) error {
	fsrc, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		log.Printf("Cannot open migrations file: %s", err)
		return err
	}
	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"postgres",
		driver)
	if err != nil {
		log.Printf("Cannot create migrate instance: %s", err)
		return err
	}
	if err := m.Steps(1); err != nil {
		log.Printf("Migration error: %s", err)
		return err
	}
	return nil
}
