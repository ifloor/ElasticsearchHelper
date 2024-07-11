package model

import (
	"ElasticsearchHelper/elasticsearch/shard_representation"
)

type WriteShardCondensedData struct {
	shard_representation.ShardCondensedData
	Timestamp string `json:"@timestamp"`
}
