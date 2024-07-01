package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	guuid "github.com/google/uuid"
	"github.com/karan2704/kube-deploy/models"
	"github.com/karan2704/kube-deploy/utils"
)

var wg sync.WaitGroup

func GenerateManifests(w http.ResponseWriter, r *http.Request){
	errorChannel := make(chan error, 10)
	var yamlStructure models.YamlConfig

	//generate a project id here and pass it as an arg to yaml generator.
	// can create a posgres global variable struct instead of one in the auth handlers.
	appId := guuid.New().String()
	trimmedId := appId[0:5]

	err := json.NewDecoder(r.Body).Decode(&yamlStructure)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = DB.Exec(models.InsertAppQuery, trimmedId, "12345")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Printf("%+v\n", yamlStructure)
	for _, app := range yamlStructure.Applications{
		go utils.YamlGenerator(app, yamlStructure.Namespace, errorChannel, trimmedId, &wg, DB)
		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for err := range errorChannel {
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
