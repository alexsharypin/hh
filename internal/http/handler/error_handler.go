package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexsharypin/hh/internal/common"

	"go.uber.org/zap"
)

func HandleError(w http.ResponseWriter, logger *zap.Logger, err error) {
	var resp common.ErrorResponse
	var status int

	if appErr, ok := err.(*common.AppError); ok {
		resp = common.ErrorResponse{
			Status:  appErr.Status,
			Message: appErr.Message,
			Errors:  appErr.Errors,
		}
		status = appErr.Status
	} else {
		resp = common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
		status = http.StatusInternalServerError

		logger.Error("unhandled error", zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
