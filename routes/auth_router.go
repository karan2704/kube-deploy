package routes

import (
	"github.com/gorilla/mux"

	"github.com/karan2704/kube-deploy/services"
)

func AuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", services.LoginHandler).Methods("POST")
}