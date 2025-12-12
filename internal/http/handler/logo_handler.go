package handler

import (
	"net/http"

	"github.com/alexsharypin/hh/internal/lib"
	"github.com/alexsharypin/hh/internal/service"
	"go.uber.org/zap"
)

type LogoHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type logoHandler struct {
	service service.LogoService
	logger  *zap.Logger
}

func NewLogoHandler(service service.LogoService, logger *zap.Logger) LogoHandler {
	return &logoHandler{
		service: service,
		logger:  logger,
	}
}

func (h *logoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := lib.ParseIDFromURL(r)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	result, err := h.service.Upload(ctx, id)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	HandleSuccess(w, http.StatusOK, result)
}

func (h *logoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := lib.ParseIDFromURL(r)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		HandleError(w, h.logger, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
