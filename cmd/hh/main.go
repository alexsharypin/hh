package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/alexsharypin/hh/internal/http/handler"
	"github.com/alexsharypin/hh/internal/repo"
	"github.com/alexsharypin/hh/internal/service"
	"github.com/go-chi/chi"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger := zap.NewExample()
	defer logger.Sync()

	pool, err := createDBPool(ctx)

	if err != nil {
		logger.Error("failed to create pg pool", zap.Error(err))
		os.Exit(1)
	}

	defer pool.Close()

	repo := repo.NewCompanyRepo(pool)

	service := service.NewCompanyService(repo)

	h := handler.NewCompanyHandler(service, logger)

	r := chi.NewRouter()

	r.Get("/companies", h.List)
	r.Post("/companies", h.Create)

	r.Route("/companies/{id}", func(r chi.Router) {
		r.Get("/", h.GetByID)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})

	logger.Debug("server started :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		logger.Error("failed to start server", zap.Error(err))
		os.Exit(1)
	}
}

func createDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return pool, nil
}
