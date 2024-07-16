package shard_representation

import (
	ilmmodel "ElasticsearchHelper/elasticsearch/ilm/model"
	shardmodel "ElasticsearchHelper/elasticsearch/sharding/model"
)

type ShardCondensedData struct {
	RunId string                  `json:"run_id"`
	Index *ilmmodel.IlmIndex      `json:"index"`
	Shard shardmodel.ElasticShard `json:"shard"`
}
