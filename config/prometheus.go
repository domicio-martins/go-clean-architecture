package config

type Prometheus struct {
	Name     string `envconfig:"PROMETHEUS_NAME" default:"template_manager_gin"`
	Path     string `envconfig:"PROMETHEUS_PATH" default:"/metrics"`
}
