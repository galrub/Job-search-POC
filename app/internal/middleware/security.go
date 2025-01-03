package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/galrub/go/jobSearch/internal/logger"
	"github.com/galrub/go/jobSearch/internal/utils"
)

var defaultTokenLookup = "header: Authorization"

func JWtMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			c, err := r.Cookie("Authorization")
			if err != nil {
				logger.LOG.Debug().Msg("cookie not found")
			} else {
				tokenString = c.Value
				http.SetCookie(w, c)
			}
		} else {
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		}
		if len(tokenString) == 0 {
			utils.RenderFragment("loginForm.html", "", &w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims, err := verifyToken(tokenString)
		if err != nil {
			var email string
			if claims != nil {
				email = claims.(jwt.MapClaims)["identity"].(string)
			} else {
				email = ""
			}
			logger.LOG.Debug().Err(err).Msg("failed to get claim from token string")
			utils.RenderFragment("loginForm.html", email, &w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		email := claims.(jwt.MapClaims)["identity"].(string)
		r.Header.Set("email", email)
		next.ServeHTTP(w, r)
	})
}
