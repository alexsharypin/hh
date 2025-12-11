package service

import (
	"context"

	"github.com/alexsharypin/hh/internal/entity"
	"github.com/alexsharypin/hh/internal/repo"
	"github.com/google/uuid"
)

type CompanyService struct {
	repo *repo.CompanyRepo
}

func NewCompanyService(repo *repo.CompanyRepo) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) Create(ctx context.Context, input *entity.CreateCompanyInput) (*entity.Company, error) {
	c, err := entity.NewCompany(input)

	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, c)
}

func (s *CompanyService) Update(ctx context.Context, id uuid.UUID, input *entity.UpdateCompanyInput) (*entity.Company, error) {
	c, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if err := c.Update(input); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, c)
}

func (s *CompanyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *CompanyService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CompanyService) GetAll(ctx context.Context) ([]*entity.Company, error) {
	return s.repo.GetAll(ctx)
}
