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
	
	//databaseObj.Db.Exec("CREATE TABLE IF NOT EXISTS Users (userId VARCHAR(20) PRIMARY KEY, username VARCHAR(20), password VARCHAR(20))")
	databaseObj.Db.Exec("CREATE TABLE IF NOT EXISTS Apps (appId VARCHAR(20) PRIMARY KEY, owner VARCHAR(20), FOREIGN KEY (owner) references Users(userId))")
	databaseObj.Db.Exec("CREATE TABLE IF NOT EXISTS Files (fileName VARCHAR(25) NOT NULL, content TEXT NOT NULL, appId VARCHAR(20), FOREIGN KEY (appId) references Apps(appId))")
	if err := http.ListenAndServe(":8080", router); err != nil{
		log.Fatalf("Error starting the server")
	}
}
