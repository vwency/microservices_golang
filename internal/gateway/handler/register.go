package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vwency/microservices_golang/internal/gateway/service"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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
			http.Error(w, fmt.Sprintf("failed to check user existence: %v", err), http.StatusInternalServerError)
			return
		}

		if resp.Found {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User successfully registered",
		})
	}
}
