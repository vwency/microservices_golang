package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vwency/microservices_golang/internal/gateway/service"
)

type ValidateRequest struct {
	AccessToken string `json:"access_token"`
}

type ValidateResponse struct {
	Valid     bool     `json:"valid"`
	UserID    string   `json:"user_id"`
	Roles     []string `json:"roles"`
	ExpiresAt int64    `json:"expires_at"`
}

func ValidateHandler(authService *service.AuthServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accessToken string

		var validateReq ValidateRequest
		if err := json.NewDecoder(r.Body).Decode(&validateReq); err == nil && validateReq.AccessToken != "" {
			accessToken = validateReq.AccessToken
		} else {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				accessToken = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if accessToken == "" {
			http.Error(w, "access token not provided", http.StatusBadRequest)
			return
		}

		resp, err := authService.Validate(r.Context(), accessToken)
		if err != nil {
			http.Error(w, "Failed to validate token", http.StatusUnauthorized)
			return
		}

		response := ValidateResponse{
			Valid:     resp.Valid,
			UserID:    resp.UserId,
			Roles:     resp.Roles,
			ExpiresAt: resp.ExpiresAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
