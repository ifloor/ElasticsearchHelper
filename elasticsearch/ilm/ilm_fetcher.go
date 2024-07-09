package ilm

import (
	"ElasticsearchHelper/elasticsearch/configs"
	"ElasticsearchHelper/elasticsearch/ilm/model"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const CatIlmPath = "/*/_ilm/explain"

func FetchIlmInfo(client *resty.Client, configs *configs.ElasticConfigs) (map[string]model.IlmIndex, error) {
	url := fmt.Sprintf("%v%v", configs.HostAndPort, CatIlmPath)

	resp, err := client.R().
		EnableTrace().
		Get(url)
	if err != nil {
		log.Error().Msgf("Error when fetching shards: %v", err)
		return map[string]model.IlmIndex{}, err
	}

	if resp.StatusCode() != 200 {
		return map[string]model.IlmIndex{}, fmt.Errorf("error when fetching ilm info: %v. Got response: %v", resp.String(), resp.StatusCode())
	}

	log.Info().Msgf("Fetch shards http Status: %v", resp.Status())

	return parseResponse(resp.String())
}

func parseResponse(body string) (map[string]model.IlmIndex, error) {
	var response model.IlmResponse

	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing string to IlmIndex: %v", err)
	}

	return response.Indices, nil
}
