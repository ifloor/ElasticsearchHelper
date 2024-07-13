package main

import (
	"ElasticsearchHelper/elasticsearch"
	"ElasticsearchHelper/elasticsearch/doc"
	"ElasticsearchHelper/elasticsearch/shard_representation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
	"strings"
)

func main() {
	setLogLevel()

	client := elasticsearch.NewElasticClient()
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

	// Prepare the data
	// Work the shards
	mappedShards := make(map[string][]*shard_representation.ShardCondensedData)
	for _, shard := range shards {
		index := shard.Index
		condensedData := shard_representation.ShardCondensedData{
			Shard: shard,
		}
		var shardsForIndex []*shard_representation.ShardCondensedData

		if _, ok := mappedShards[index]; ok {
			shardsForIndex = mappedShards[index]
			shardsForIndex = append(shardsForIndex, &condensedData)
		} else {
			shardsForIndex = []*shard_representation.ShardCondensedData{&condensedData}
		}
		mappedShards[index] = shardsForIndex
	}
	// Work the ilm indices
	for _, ilm := range ilmInfo {
		if condensedList, ok := mappedShards[ilm.IndexName]; ok {
			for _, condensed := range condensedList {
				condensed.Index = ilm
			}
		} else {
			log.Warn().Msgf("ILM Index %v does not have any shards found previously...", ilm.IndexName)
		}
	}

	// Prepare the list
	var condensedShards = make([]*shard_representation.ShardCondensedData, 0)
	for _, shards := range mappedShards {
		condensedShards = append(condensedShards, shards...)
	}

	// Sending to Elasticsearch
	writeShardsData(client, condensedShards)
}

func writeShardsData(client *elasticsearch.ElasticClient, condensedShards []*shard_representation.ShardCondensedData) {
	config := WritingElasticsearchConfig{}
	if err := envconfig.Init(&config); err != nil {
		log.Error().Msgf("Error when reading Writing elasticsearch configs: %v", err)
		panic(err)
	}
	for _, condensedShard := range condensedShards {
		err := doc.WriteDocument(client.Client, client.Configs, config.IndexName, condensedShard)
		if err != nil {
			log.Error().Msgf("Error when writing document: %v", err)
		}
	}
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

type WritingElasticsearchConfig struct {
	IndexName string `envconfig:"ELASTIC_WRITE_INDEX_NAME"`
}
