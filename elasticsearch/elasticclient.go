package elasticsearch

import (
	configs2 "ElasticsearchHelper/elasticsearch/configs"
	"ElasticsearchHelper/elasticsearch/ilm"
	"ElasticsearchHelper/elasticsearch/ilm/model"
	"ElasticsearchHelper/elasticsearch/sharding"
	modelsharding "ElasticsearchHelper/elasticsearch/sharding/model"
	"github.com/go-resty/resty/v2"
)

type ElasticClient struct {
	Client  *resty.Client
	Configs *configs2.ElasticConfigs
}

func NewElasticClient() *ElasticClient {
	configs := configs2.GetElasticConfigs()
	client := resty.New()
	client = client.SetBasicAuth(configs.AuthUsername, configs.AuthUserPassword)
	return &ElasticClient{
		Client:  client,
		Configs: configs,
	}
}

func (ec *ElasticClient) FetchShards() ([]modelsharding.ElasticShard, error) {
	return sharding.FetchShards(ec.Client, ec.Configs)
}

func (ec *ElasticClient) FetchIlmInfo() (map[string]model.IlmIndex, error) {
	return ilm.FetchIlmInfo(ec.Client, ec.Configs)
}
