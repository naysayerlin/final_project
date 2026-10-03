package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

const jwtSecret = "final_project_secret_key"

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, signInResponse{Error: "Unsupported method"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, signInResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	storedPassword := appConfig.Password

	if req.Password != storedPassword {
		writeJSON(w, http.StatusUnauthorized, signInResponse{Error: "invalid password"})
		return
	}

	hash := passwordHash(storedPassword)

	claims := jwt.MapClaims{
		"hash": hash,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, signInResponse{Error: "failed to sign token: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, signInResponse{Token: signed})
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum)
}
