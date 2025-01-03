package utils

import (
	"encoding/json"
	"net/http"

	"github.com/galrub/go/jobSearch/internal/logger"
	"github.com/galrub/go/jobSearch/internal/model"
	"github.com/galrub/go/jobSearch/internal/services"
)

func LoginVerification(w http.ResponseWriter, r *http.Request) (string, bool) {
	var loginData model.LoginData

	requestCtx := r.Header.Get("content-type")
	if requestCtx != "application/json" {
		email := r.PostFormValue("email")
		pw := r.PostFormValue("password")
		if email == "" || pw == "" {
			logger.LOG.Error().Msg("one or more form fields missing")
			http.Error(w, "some form field are missing", http.StatusBadRequest)
			return "", false
		}
		loginData = model.LoginData{Email: email, Passw: pw}
	} else {
		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			logger.LOG.Err(err).Msg("failed to parse login data")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return "", false
		}
	}

	res, err := services.LoginUserByEmail(loginData.Email, loginData.Passw)
	if err != nil {
		logger.LOG.Err(err).Msg("failed to login user")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", false
	}

	return loginData.Email, res
}
