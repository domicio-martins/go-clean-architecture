package config

type Database struct {
	MongoAddress       string `envconfig:"MONGO_ADDRESS" default:"mongodb://localhost:27018"`
	MongoTimeout       int    `envconfig:"MONGO_TIMEOUT" default:"10"`
	MongoDatabase      string `envconfig:"MONGO_DATABASE" default:"template_manager"`
	Client
}
