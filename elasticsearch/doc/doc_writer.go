package doc

import (
	"ElasticsearchHelper/elasticsearch/configs"
	"ElasticsearchHelper/elasticsearch/doc/model"
	"ElasticsearchHelper/elasticsearch/shard_representation"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"time"
)

func WriteDocument(
	client *resty.Client,
	configs *configs.ElasticConfigs,
	indexName string,
	object *shard_representation.ShardCondensedData,
) error {
	outObject := model.WriteShardCondensedData{
		ShardCondensedData: *object,
		Timestamp:          time.Now().Format(time.RFC3339),
	}

	url := fmt.Sprintf("%v/%v/_doc", configs.HostAndPort, indexName)

	resp, err := client.R().
		SetBody(outObject).
		EnableTrace().
		Post(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("error when writing doc info: %v. Got response: %v", resp.String(), resp.StatusCode())
	}

	log.Info().Msgf("Wrote _doc successfully. Http Status: %v", resp.Status())

	return nil
}
