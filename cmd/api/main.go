package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/EduardoMark/my-finance-api/pkg/config"
	"github.com/EduardoMark/my-finance-api/pkg/database"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbInstance, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbInstance.Close()

	srv := http.Server{
		Addr:              ":3000",
		Handler:           nil,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Minute * 1,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
