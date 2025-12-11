package service

import (
	"context"

	"github.com/alexsharypin/hh/internal/repo"
)

type LogoService interface {
	Upload(ctx context.Context) error
	Delete(ctx context.Context) error
}

type logoService struct {
	repo *repo.CompanyRepo
}

func NewLogoService(repo *repo.CompanyRepo) LogoService {
	return &logoService{repo: repo}
}

func (s *logoService) Upload(ctx context.Context) error {
	return ctx.Err()
}

func (s *logoService) Delete(ctx context.Context) error {
	return ctx.Err()
}
