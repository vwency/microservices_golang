package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vwency/microservices_golang/internal/gateway/service"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func RegisterHandler(authService *service.AuthServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerReq RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := authService.GetUser(r.Context(), registerReq.Username, registerReq.Email)
		if err != nil {
			http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
			return
		}

		if resp.Found {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		authResp, err := authService.Register(r.Context(), registerReq.Username, registerReq.Password, registerReq.Email)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		response := RegisterResponse{
			AccessToken:  authResp.AccessToken,
			RefreshToken: authResp.RefreshToken,
			ExpiresAt:    authResp.ExpiresAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
