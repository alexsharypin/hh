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
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	minio, err := createMinio()
	if err != nil {
		logger.Error("failed to create minio client", zap.Error(err))
		os.Exit(1)
	}

	companyRepo := repo.NewCompanyRepo(pool)

	companyService := service.NewCompanyService(companyRepo)

	logoService := service.NewLogoService(minio, companyRepo)

	companyHandler := handler.NewCompanyHandler(companyService, logger)

	logoHandler := handler.NewLogoHandler(logoService, logger)

	r := chi.NewRouter()

	r.Get("/companies", companyHandler.List)
	r.Post("/companies", companyHandler.Create)

	r.Route("/companies/{id}", func(r chi.Router) {
		r.Get("/", companyHandler.GetByID)
		r.Put("/", companyHandler.Update)
		r.Delete("/", companyHandler.Delete)

		r.Post("/logo", logoHandler.Upload)
		r.Delete("/logo", logoHandler.Delete)
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

func createMinio() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
