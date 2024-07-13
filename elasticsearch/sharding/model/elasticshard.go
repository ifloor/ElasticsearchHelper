package model

import "strconv"

type ElasticShard struct {
	Index            string
	Shard            string
	PriRep           string
	State            string
	Node             string
	UnassignedReason string
	Docs             string
	Store            int64
}

func (es *ElasticShard) ToString() string {
	return "Index=" + es.Index + ", Shard=" + es.Shard + ", PriRep=" + es.PriRep + ", State=" + es.State +
		", Node=" + es.Node + ", UnassignedReason=" + es.UnassignedReason + ", Docs=" + es.Docs + ", Store=" + strconv.FormatInt(es.Store, 10)
}
