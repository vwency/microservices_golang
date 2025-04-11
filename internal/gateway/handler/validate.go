package handler

import (
	"encoding/json"
	"net/http"

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
		var validateReq ValidateRequest
		if err := json.NewDecoder(r.Body).Decode(&validateReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := authService.Validate(r.Context(), validateReq.AccessToken)
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
