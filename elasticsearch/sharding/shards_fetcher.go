package sharding

import (
	"ElasticsearchHelper/elasticsearch/configs"
	model2 "ElasticsearchHelper/elasticsearch/sharding/model"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

const CatShardsPath = "/_cat/shards"
const CatShardsParams = "?v=true&h=index,shard,prirep,state,node,unassigned.reason,docs,store,dataset.size"

type ShardsFetcher struct{}

type TempToken struct {
	Token         string
	FirstPosition int
}

func FetchShards(client *resty.Client, configs *configs.ElasticConfigs) ([]model2.ElasticShard, error) {
	url := fmt.Sprintf("%v%v%v", configs.HostAndPort, CatShardsPath, CatShardsParams)

	resp, err := client.R().
		EnableTrace().
		Get(url)
	if err != nil {
		log.Error().Msgf("Error when fetching shards: %v", err)
		return []model2.ElasticShard{}, err
	}

	if resp.StatusCode() != 200 {
		return []model2.ElasticShard{}, fmt.Errorf("error when fetching shards: %v. Got response: %v", resp.String(), resp.StatusCode())
	}

	log.Info().Msgf("Fetch shards http Status: %v", resp.Status())

	return parseResponse(resp.String())
}

func parseResponse(response string) ([]model2.ElasticShard, error) {
	var shards []model2.ElasticShard

	var headers []model2.ShardHeaderPositioned
	pieces := strings.Split(response, "\n")
	for index, piece := range pieces {
		if index == 0 {
			headers = parseHeaderLine(piece)
			for _, header := range headers {
				log.Debug().Msgf("Header: %v", header)
			}
			continue
		}

		shardInfo, err := parseLine(piece, headers)
		if err != nil {
			return []model2.ElasticShard{}, err
		}
		shards = append(shards, shardInfo)
		log.Debug().Msgf("Shard: %v", shardInfo.ToString())
	}

	return shards, nil
}
