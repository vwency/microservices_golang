package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vwency/microservices_golang/internal/gateway/service"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func RefreshHandler(authService *service.AuthServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshReq RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := authService.Refresh(r.Context(), refreshReq.RefreshToken)
		if err != nil {
			http.Error(w, "Failed to refresh tokens", http.StatusUnauthorized)
			return
		}

		response := RefreshResponse{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    resp.ExpiresAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
