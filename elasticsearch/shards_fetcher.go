package elasticsearch

import (
	"ElasticsearchHelper/elasticsearch/model"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

const CatShardsPath = "/_cat/shards"
const CatShardsParams = "?v=true&h=index,shard,prirep,state,node,unassigned.reason,docs,store,dataset.size"

type ShardsFetcher struct {
	client  *resty.Client
	configs *ElasticConfigs
}

type TempToken struct {
	Token         string
	FirstPosition int
}

func NewShardsFetcher() *ShardsFetcher {
	configs := GetElasticConfigs()
	client := resty.New()
	client = client.SetBasicAuth(configs.AuthUsername, configs.AuthUserPassword)
	return &ShardsFetcher{
		client:  client,
		configs: configs,
	}
}

func (s *ShardsFetcher) FetchShards() ([]model.ElasticShard, error) {
	url := fmt.Sprintf("%v%v%v", s.configs.HostAndPort, CatShardsPath, CatShardsParams)

	resp, err := s.client.R().
		EnableTrace().
		Get(url)
	if err != nil {
		log.Error().Msgf("Error when fetching shards: %v", err)
		return []model.ElasticShard{}, err
	}

	if resp.StatusCode() != 200 {
		return []model.ElasticShard{}, fmt.Errorf("error when fetching shards: %v. Got response: %v", resp.String(), resp.StatusCode())
	}

	log.Info().Msgf("Status: %v", resp.Status())

	return parseResponse(resp.String())
}

func parseResponse(response string) ([]model.ElasticShard, error) {
	var shards []model.ElasticShard

	var headers []model.ShardHeaderPositioned
	pieces := strings.Split(response, "\n")
	for index, piece := range pieces {
		fmt.Println(piece)
		if index == 0 {
			headers = parseHeaderLine(piece)
			for _, header := range headers {
				log.Debug().Msgf("Header: %v", header)
			}
			continue
		}

		shardInfo, err := parseLine(piece, headers)
		if err != nil {
			return []model.ElasticShard{}, err
		}
		shards = append(shards, shardInfo)
		log.Debug().Msgf("Shard: %v", shardInfo.ToString())
	}

	return shards, nil
}
