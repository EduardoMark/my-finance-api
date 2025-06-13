package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/EduardoMark/my-finance-api/internal/api"
	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/pkg/config"
	"github.com/EduardoMark/my-finance-api/pkg/database"
	"github.com/EduardoMark/my-finance-api/pkg/token"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	db := db.New(dbPool)

	token := token.NewTokenManager(*cfg)

	apiInstance := api.NewApi(cfg, db, token)
	apiInstance.SetupApi()
	apiInstance.BindRoutes()

	srv := http.Server{
		Addr:              ":3000",
		Handler:           apiInstance.Router,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Minute * 1,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
