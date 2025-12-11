package handler

import (
	"net/http"

	"github.com/alexsharypin/hh/internal/service"
	"go.uber.org/zap"
)

type LogoHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type logoHandler struct {
	service *service.CompanyService
	logger  *zap.Logger
}

func NewLogoHandler(service *service.CompanyService, logger *zap.Logger) LogoHandler {
	return &logoHandler{
		service: service,
		logger:  logger,
	}
}

func (h *logoHandler) Upload(w http.ResponseWriter, r *http.Request) {}

func (h *logoHandler) Delete(w http.ResponseWriter, r *http.Request) {}
