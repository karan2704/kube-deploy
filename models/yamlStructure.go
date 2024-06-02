package models

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}

type Container struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
	Ports []Port `yaml:"ports"`
	Env   []EnvVar `yaml:"env,omitempty"`
}

type Port struct {
	ContainerPort int `yaml:"containerPort"`
}

type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type PodSpec struct {
	Containers []Container `yaml:"containers"`
}

type PodTemplateSpec struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     PodSpec `yaml:"spec"`
}

type Selector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type DeploymentSpec struct {
	Replicas int              `yaml:"replicas"`
	Selector Selector         `yaml:"selector"`
	Template PodTemplateSpec  `yaml:"template"`
}

type Deployment struct {
	ApiVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   Metadata      `yaml:"metadata"`
	Spec       DeploymentSpec `yaml:"spec"`
}

type ServicePort struct {
	Protocol   string `yaml:"protocol"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
}

type ServiceSpec struct {
	Selector map[string]string `yaml:"selector"`
	Ports    []ServicePort     `yaml:"ports"`
	Type     string            `yaml:"type"`
}

type Service struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       ServiceSpec `yaml:"spec"`
}

type ConfigMap struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}
