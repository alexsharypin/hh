package lib

import (
	"net/http"

	"github.com/alexsharypin/hh/internal/common"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func ParseIDFromURL(r *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return uuid.Nil, common.NewInvalidRequestParams()
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, common.NewInvalidRequestParams()
	}

	return id, nil
}
