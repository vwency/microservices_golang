package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vwency/microservices_golang/internal/gateway/service"
)

type LogoutRequest struct {
	Username string `json:"username"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func LogoutHandler(authService *service.AuthServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logoutReq LogoutRequest
		if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		resp, err := authService.Logout(r.Context(), logoutReq.Username, accessToken)
		if err != nil {
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
			return
		}

		response := LogoutResponse{
			Success: resp.Success,
			Message: resp.Message,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
