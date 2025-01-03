package handlers

import (
	"net/http"

	"github.com/galrub/go/jobSearch/internal/logger"
)

func JobsByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	if email == "" {
		logger.LOG.Error().Msg("email not in request headers")
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
