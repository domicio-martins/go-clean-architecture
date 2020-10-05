package config

type General struct {
	Debug              int    `envconfig:"DEBUG" default:"0"`
	AppEnv             string `envconfig:"APP_ENV" default:"prod"`
	Host               string `default:"http://localhost:9000"`
	HttpAddress        string `envconfig:"HTTP_ADDRESS" default:"0.0.0.0:9000"`
}
