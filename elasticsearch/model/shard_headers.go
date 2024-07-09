package model

type ShardHeader string

const (
	Index            ShardHeader = "index"
	Shard            ShardHeader = "shard"
	PriRep           ShardHeader = "prirep"
	State            ShardHeader = "state"
	Node             ShardHeader = "node"
	UnassignedReason ShardHeader = "unassigned.reason"
	Docs             ShardHeader = "docs"
	Store            ShardHeader = "store"
)

type ShardHeaderPositioned struct {
	Header        ShardHeader
	FirstPosition int
	LastPosition  int
}
