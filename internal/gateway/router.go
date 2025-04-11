package gateway

import (
	"github.com/gorilla/mux"
	"github.com/vwency/microservices_golang/internal/gateway/handler"
	"github.com/vwency/microservices_golang/internal/gateway/service"
)

func InitializeRouter(authService *service.AuthServiceClient) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", handler.RegisterHandler(authService)).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler(authService)).Methods("POST")

	return r
}
