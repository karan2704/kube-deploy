package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	guuid "github.com/google/uuid"
	"github.com/karan2704/kube-deploy/models"
	"gopkg.in/yaml.v2"
)



func YamlGenerator(fields models.Application, namespace string, errorChannel chan<-error, wg *sync.WaitGroup, db *sql.DB) () {
	defer wg.Done()

	secrets := unsetSecrets(fields)

	deployment := models.Deployment{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: models.Metadata{
			Name:      fields.Name + "-deployment",
			Namespace: namespace,
		},
		Spec: models.DeploymentSpec{
			Replicas: 2,
			Selector: models.Selector{
				MatchLabels: map[string]string{"app": fields.Name},
			},
			Template: models.PodTemplateSpec{
				Metadata: models.Metadata{
					Name: fields.Name,
					Namespace: namespace,
					Labels: map[string]string{"app": fields.Name},
				},
				Spec: models.PodSpec{
					Containers: []models.Container{
						{
							Name:  fields.Name + "-container",
							Image: fields.Image,
							Ports: []models.Port{{ContainerPort: fields.Ports[0]}},
						},
					},
				},
			},
		},
	}

	service := models.Service{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: models.Metadata{
			Name:      fields.Name + "-service",
			Namespace: namespace,
		},
		Spec: models.ServiceSpec{
			Selector: map[string]string{"app": fields.Name},
			Ports: []models.ServicePort{
				{Protocol: "TCP", Port: fields.Ports[0], TargetPort: fields.Ports[0]},
			},
			Type: "LoadBalancer",
		},
	}

	configMap := models.ConfigMap{
		ApiVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: models.Metadata{
			Name:      fields.Name + "-config",
			Namespace: namespace,
		},
		Data: *secrets,
	}

	deploymentYAML, err := yaml.Marshal(deployment)
	if err != nil {
		errorChannel <- fmt.Errorf("error marshalling yaml")
		return
	}
	serviceYAML, err := yaml.Marshal(service)
	if err != nil {
		errorChannel <- fmt.Errorf("error marshalling yaml")
		return
	}
	configMapYAML, err := yaml.Marshal(configMap)
	if err != nil {
		errorChannel <- fmt.Errorf("error marshalling yaml")
		return
	}

	// Combine all YAMLs into one
	fullYAML := string(deploymentYAML) + "---\n" + string(serviceYAML) + "---\n" + string(configMapYAML)

	// Write to a file
	directory := "files"
	fileName := fields.Name + "-manifest.yaml"
	filePath := filepath.Join(directory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		errorChannel <- fmt.Errorf("error creating yaml file %s", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fullYAML)
	if err != nil {
		errorChannel <- fmt.Errorf("error writing to yaml file %s", err)
		return
	}

	fmt.Printf("%s manifest.yaml created successfully \n", fields.Name)

	_, err = db.Exec(models.InsertManifestFiles, fields.Name, fullYAML)

	return
}

func unsetSecrets(fields models.Application) (*map[string]string){
	secrets := make(map[string]string);
	secrets["appName"] = fields.Name
	envVars := reflect.ValueOf(fields.Secrets)
	types := envVars.Type()

	for i:=0; i < envVars.NumField(); i++ {
		secrets[types.Field(i).Name] = envVars.Field(i).String()
	}
	return &secrets
}

func getUUID()(string){
	return guuid.New()
}