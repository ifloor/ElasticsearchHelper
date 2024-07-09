package main

import (
	"ElasticsearchHelper/elasticsearch"
	"fmt"
	"github.com/rs/zerolog/log"
)

func main() {
	//log.Le

	fmt.Println("Hello, World!")
	shards, err := elasticsearch.NewShardsFetcher().FetchShards()
	if err != nil {
		log.Error().Msgf("Error when fetching shards: %v", err)
	}
	log.Info().Msgf("Found shards: %v", len(shards))

}
