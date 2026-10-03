package api

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := appConfig.Password
		if len(pass) > 0 {
			var tokenStr string
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenStr = cookie.Value
			}
			valid := validateToken(tokenStr, pass)
			if !valid {
				writeError(w, http.StatusUnauthorized, "Authentification required")
				return
			}
		}
		next(w, r)
	}
}

func validateToken(tokenStr string, currentPassword string) bool {
	if tokenStr == "" {
		return false
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	hash, ok := claims["hash"].(string)
	if !ok {
		return false
	}
	return hash == passwordHash(currentPassword)
}
