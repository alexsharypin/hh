package service

import (
	"context"
	"time"

	"github.com/alexsharypin/hh/internal/repo"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const (
	LogosBucket  = "logos"
	LogosExpires = time.Hour
	LogosMinSize = 1 * 1024 * 256
	LogosMaxSize = 1 * 1024 * 1024
)

type UploadResult struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

type LogoService interface {
	Upload(ctx context.Context, id uuid.UUID) (*UploadResult, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type logoService struct {
	repo  repo.CompanyRepo
	minio *minio.Client
}

func NewLogoService(minio *minio.Client, repo repo.CompanyRepo) LogoService {
	return &logoService{minio: minio, repo: repo}
}

func (s *logoService) Upload(ctx context.Context, id uuid.UUID) (*UploadResult, error) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	policy := minio.NewPostPolicy()
	policy.SetBucket(LogosBucket)
	policy.SetKey(id.String())
	policy.SetExpires(time.Now().UTC().Add(LogosExpires))
	policy.SetContentLengthRange(LogosMinSize, LogosMaxSize)

	url, data, err := s.minio.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	for k, v := range data {
		fields[k] = v
	}

	return &UploadResult{
		URL:    url.String(),
		Fields: fields,
	}, nil
}

func (s *logoService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.minio.RemoveObject(ctx, LogosBucket, id.String(), minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
