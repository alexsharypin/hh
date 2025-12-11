package entity

import (
	"net/url"
	"strings"
	"time"

	"github.com/alexsharypin/hh/internal/common"
	"github.com/google/uuid"
)

type Company struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Website     *string   `json:"website"`
	LogoURL     *string   `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCompanyInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Website     *string `json:"website,omitempty"`
}

type UpdateCompanyInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Website     *string `json:"website,omitempty"`
}

func NewCompany(input *CreateCompanyInput) (*Company, error) {
	description := ""

	if input.Description != nil {
		description = strings.TrimSpace(*input.Description)
	}

	c := &Company{
		Title:       strings.TrimSpace(input.Title),
		Description: description,
		Website:     input.Website,
	}

	c.beforeSave()

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Company) Update(input *UpdateCompanyInput) error {
	if input.Title != nil {
		c.Title = strings.TrimSpace(*input.Title)
	}

	if input.Description != nil {
		c.Description = strings.TrimSpace(*input.Description)
	}

	if input.Website != nil {
		c.Website = input.Website
	}

	c.beforeSave()

	return c.Validate()
}

func (c *Company) beforeSave() {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	c.UpdatedAt = time.Now().UTC()
}

func (c *Company) Validate() error {
	var errs []common.ValidationError

	c.Title = strings.TrimSpace(c.Title)
	if len(c.Title) < 3 || len(c.Title) > 64 {
		errs = append(errs, common.ValidationError{
			Field:   "title",
			Message: "must be between 3 and 64 characters",
		})
	}

	c.Description = strings.TrimSpace(c.Description)
	if len(c.Description) > 255 {
		errs = append(errs, common.ValidationError{
			Field:   "description",
			Message: "must be at most 255 characters",
		})
	}

	if c.Website != nil {
		website := strings.TrimSpace(*c.Website)
		u, err := url.Parse(website)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			errs = append(errs, common.ValidationError{
				Field:   "website",
				Message: "must be a valid URL starting with http:// or https://",
			})
		}
		c.Website = &website
	}

	if len(errs) > 0 {
		return common.NewValidationError(errs)
	}

	return nil
}
