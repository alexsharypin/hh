package service

import (
	"context"

	"github.com/alexsharypin/hh/internal/entity"
	"github.com/alexsharypin/hh/internal/repo"
	"github.com/google/uuid"
)

type CompanyService interface {
	Create(ctx context.Context, input *entity.CreateCompanyInput) (*entity.Company, error)
	Update(ctx context.Context, id uuid.UUID, input *entity.UpdateCompanyInput) (*entity.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error)
	GetAll(ctx context.Context) ([]*entity.Company, error)
}

type companyService struct {
	repo repo.CompanyRepo
}

func NewCompanyService(repo repo.CompanyRepo) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) Create(ctx context.Context, input *entity.CreateCompanyInput) (*entity.Company, error) {
	c, err := entity.NewCompany(input)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, c)
}

func (s *companyService) Update(ctx context.Context, id uuid.UUID, input *entity.UpdateCompanyInput) (*entity.Company, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := c.Update(input); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, c)
}

func (s *companyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *companyService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *companyService) GetAll(ctx context.Context) ([]*entity.Company, error) {
	return s.repo.GetAll(ctx)
}
