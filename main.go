package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/karan2704/kube-deploy/config"
	"github.com/karan2704/kube-deploy/routes"
)

func main(){
	db, err := config.ConfigureDB()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", db)

	router := mux.NewRouter()

	routes.AuthRoutes(router)
	routes.KubeRoutes(router)
	
	if err := http.ListenAndServe(":8080", router); err != nil{
		log.Fatalf("Error starting the server")
	}
}
