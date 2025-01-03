package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/galrub/go/jobSearch/internal/logger"
	"github.com/galrub/go/jobSearch/internal/model"
	"github.com/galrub/go/jobSearch/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func Login(w http.ResponseWriter, r *http.Request) {
	email, res := utils.LoginVerification(w, r)
	if res == false {
		resp := model.LoginResponse{Status: "failed", Token: "", Message: "Login Failed"}
		encoder := json.NewEncoder(w)
		err := encoder.Encode(resp)
		if err != nil {
			logger.LOG.Fatal().Err(err).Msg("cannot encode json response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = email
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		logger.LOG.Fatal().Err(err).Msg("cannot encode generate jwt token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "authentication",
		Value:    t,
		Expires:  time.Now().Add(2 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)

}
