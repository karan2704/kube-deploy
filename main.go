package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/karan2704/kube-deploy/config"
	"github.com/karan2704/kube-deploy/routes"
	"github.com/karan2704/kube-deploy/services"
)


func main(){
	databaseObj, err := config.ConfigureDB()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", databaseObj)

	services.SetDB(databaseObj.Db)

	router := mux.NewRouter()

	routes.AuthRoutes(router)
	routes.KubeRoutes(router)
	
	services.AuthInit()

	if err := http.ListenAndServe(":8080", router); err != nil{
		log.Fatalf("Error starting the server")
	}
}
