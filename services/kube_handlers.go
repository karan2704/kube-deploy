package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/karan2704/kube-deploy/models"
	"github.com/karan2704/kube-deploy/utils"
)

var wg sync.WaitGroup

func GenerateManifests(w http.ResponseWriter, r *http.Request){
	errorChannel := make(chan error, 10)
	var yamlStructure models.YamlConfig

	err := json.NewDecoder(r.Body).Decode(&yamlStructure)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Printf("%+v\n", yamlStructure)
	for _, app := range yamlStructure.Applications{
		go utils.YamlGenerator(app, yamlStructure.Namespace, errorChannel, &wg)
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