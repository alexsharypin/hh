package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexsharypin/hh/internal/common"
	"github.com/alexsharypin/hh/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepo interface {
	Create(ctx context.Context, c *entity.Company) (*entity.Company, error)
	Update(ctx context.Context, c *entity.Company) (*entity.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error)
	GetAll(ctx context.Context) ([]*entity.Company, error)
}

type companyRepo struct {
	db *pgxpool.Pool
}

func NewCompanyRepo(db *pgxpool.Pool) CompanyRepo {
	return &companyRepo{db: db}
}

func (r *companyRepo) Create(ctx context.Context, c *entity.Company) (*entity.Company, error) {
	query := "INSERT INTO companies (id, title, description, website, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.Exec(ctx, query, c.ID, c.Title, c.Description, c.Website, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, common.NewValidationError([]common.ValidationError{
					{Field: "title", Message: "must be unique"},
				})
			}
		}

		return nil, err
	}

	return c, nil
}

func (r *companyRepo) Update(ctx context.Context, c *entity.Company) (*entity.Company, error) {
	query := "UPDATE companies SET title=$2, description=$3, website=$4, created_at=$5, updated_at=$6 WHERE id=$1"

	_, err := r.db.Exec(ctx, query, c.ID, c.Title, c.Description, c.Website, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewNotFoundError("Company not found")
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, common.NewValidationError([]common.ValidationError{
					{Field: "title", Message: "must be unique"},
				})
			}
		}

		return nil, err
	}

	return c, nil
}

func (r *companyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM companies WHERE id=$1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return common.NewNotFoundError("Company not found")
	}

	return nil
}

func (r *companyRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	c := &entity.Company{}

	query := "SELECT id, title, description, website, logo_url, created_at, updated_at FROM companies WHERE id=$1"

	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Title, &c.Description, &c.Website, &c.LogoURL, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewNotFoundError("Company not found")
		}

		return nil, err
	}

	return c, nil
}

func (r *companyRepo) GetAll(ctx context.Context) ([]*entity.Company, error) {
	query := "SELECT id, title, description, website, logo_url, created_at, updated_at FROM companies"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []*entity.Company{}

	for rows.Next() {
		c := &entity.Company{}

		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Website, &c.LogoURL, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, nil
}
