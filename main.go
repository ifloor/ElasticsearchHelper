package main

import (
	"ElasticsearchHelper/elasticsearch"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
	"strings"
)

func main() {
	setLogLevel()

	client := elasticsearch.NewElasticClient()
	fmt.Println("Hello, World!")
	shards, err := client.FetchShards()
	if err != nil {
		log.Error().Msgf("Error when fetching shards: %v", err)
		return
	}
	log.Info().Msgf("Found shards: %v", len(shards))

	ilmInfo, err := client.FetchIlmInfo()
	if err != nil {
		log.Error().Msgf("Error when fetching ilm info: %v", err)
		return
	}
	log.Debug().Msgf("Found ilm info: %v", ilmInfo)
	log.Info().Msgf("Found ilm infos: %v", len(ilmInfo))

}

func setLogLevel() {
	config := LogConfig{}
	if err := envconfig.Init(&config); err != nil {
		log.Warn().Msgf("Error when reading log configs: %v", err)
		config = LogConfig{
			LogLevel: "info",
		}
	}
	if strings.ToLower(config.LogLevel) == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if strings.ToLower(config.LogLevel) == "info" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if strings.ToLower(config.LogLevel) == "warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if strings.ToLower(config.LogLevel) == "error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

type LogConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL"`
}
