package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexsharypin/hh/internal/common"
	"github.com/alexsharypin/hh/internal/entity"
	"github.com/alexsharypin/hh/internal/lib"
	"github.com/alexsharypin/hh/internal/service"

	"go.uber.org/zap"
)

type CompanyHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}

type companyHandler struct {
	service service.CompanyService
	logger  *zap.Logger
}

func NewCompanyHandler(service service.CompanyService, logger *zap.Logger) CompanyHandler {
	return &companyHandler{
		service: service,
		logger:  logger,
	}
}

func (h *companyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := &entity.CreateCompanyInput{}

	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		HandleError(w, h.logger, common.NewInvalidRequestBody())
		return
	}

	result, err := h.service.Create(ctx, input)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	HandleSuccess(w, http.StatusCreated, result)
}

func (h *companyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := lib.ParseIDFromURL(r)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	input := &entity.UpdateCompanyInput{}

	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		HandleError(w, h.logger, common.NewInvalidRequestBody())
		return
	}

	result, err := h.service.Update(ctx, id, input)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	HandleSuccess(w, http.StatusOK, result)
}

func (h *companyHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *companyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.service.GetAll(ctx)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	HandleSuccess(w, http.StatusOK, result)
}

func (h *companyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := lib.ParseIDFromURL(r)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	result, err := h.service.GetByID(ctx, id)
	if err != nil {
		HandleError(w, h.logger, err)
		return
	}

	HandleSuccess(w, http.StatusOK, result)
}
