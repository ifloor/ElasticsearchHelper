package elasticsearch

import (
	configs2 "ElasticsearchHelper/elasticsearch/configs"
	"ElasticsearchHelper/elasticsearch/ilm"
	"ElasticsearchHelper/elasticsearch/ilm/model"
	"ElasticsearchHelper/elasticsearch/sharding"
	model2 "ElasticsearchHelper/elasticsearch/sharding/model"
	"github.com/go-resty/resty/v2"
)

type ElasticClient struct {
	client  *resty.Client
	configs *configs2.ElasticConfigs
}

func NewElasticClient() *ElasticClient {
	configs := configs2.GetElasticConfigs()
	client := resty.New()
	client = client.SetBasicAuth(configs.AuthUsername, configs.AuthUserPassword)
	return &ElasticClient{
		client:  client,
		configs: configs,
	}
}

func (ec *ElasticClient) FetchShards() ([]model2.ElasticShard, error) {
	return sharding.FetchShards(ec.client, ec.configs)
}

func (ec *ElasticClient) FetchIlmInfo() (map[string]model.IlmIndex, error) {
	return ilm.FetchIlmInfo(ec.client, ec.configs)
}
