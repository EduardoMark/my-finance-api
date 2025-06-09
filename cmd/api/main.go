package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/EduardoMark/my-finance-api/internal/db"
	"github.com/EduardoMark/my-finance-api/internal/user"
	"github.com/EduardoMark/my-finance-api/pkg/config"
	"github.com/EduardoMark/my-finance-api/pkg/database"
	"github.com/EduardoMark/my-finance-api/pkg/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc, token)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
	})

	srv := http.Server{
		Addr:              ":3000",
		Handler:           r,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Minute * 1,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
