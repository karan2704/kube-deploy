package models

type Secrets struct {
	Env1 string `json:"env1"`
	Env2 string `json:"env2"`
}

// Define the struct for each application
type Application struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Ports   []int    `json:"ports"`
	Secrets Secrets  `json:"secrets"`
}

// Define the struct for the main JSON structure
type YamlConfig struct {
	Namespace    string        `json:"namespace"`
	Applications []Application `json:"applications"`
}