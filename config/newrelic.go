package config

type NewRelic struct {
	LicenseKey string `envconfig:"NEW_RELIC_LICENSE_KEY"`
}
