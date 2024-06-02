package routes

import (
	"github.com/gorilla/mux"

	"github.com/karan2704/kube-deploy/services"
)

func KubeRoutes(router *mux.Router){
	router.HandleFunc("/manifest/generate", services.GenerateManifests).Methods("POST")
	// router.HandleFunc("/manifest/update", services.UpdateManifests())
	// router.HandleFunc("/kube/deploy", services.Deploy())
	// router.HandleFunc("/kube/restart", services.Restart())
	// router.HandleFunc("/kube/view", services.ViewResources())
	// router.HandleFunc("/kube/logs", services.ViewLogs())
}