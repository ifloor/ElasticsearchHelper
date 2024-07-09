package elasticsearch

import (
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
)

type ElasticConfigs struct {
	HostAndPort      string `envconfig:"ELASTICSEARCH_HOST_AND_PORT"`
	AuthUsername     string `envconfig:"ELASTICSEARCH_USERNAME"`
	AuthUserPassword string `envconfig:"ELASTICSEARCH_PASSWORD"`
}

func GetElasticConfigs() *ElasticConfigs {
	config := ElasticConfigs{}
	if err := envconfig.Init(&config); err != nil {
		log.Error().Msgf("Error when reading elasticsearch configs: %v", err)
		panic(err)
	}

	if config.HostAndPort == "" {
		panic("environment variable ELASTICSEARCH_HOST_AND_PORT is required")
	}

	if config.AuthUsername == "" {
		panic("environment variable ELASTICSEARCH_USERNAME is required")
	}

	if config.AuthUserPassword == "" {
		panic("environment variable ELASTICSEARCH_PASSWORD is required")
	}

	return &config
}
