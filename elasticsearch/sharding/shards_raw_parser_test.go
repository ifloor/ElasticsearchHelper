package sharding

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type TestCase struct {
	TestName         string
	Input            string
	Index            string
	Shard            string
	PriRep           string
	State            string
	Node             string
	UnassignedReason string
	Docs             string
	Store            int64
}

var TestCases = []TestCase{
	{
		TestName: "Normal",
		Input: "index                                                                    shard prirep state      node                                                        unassigned.reason     docs   store\n" +
			".apm-agent-configuration                                                 0     r      STARTED    aragorn                                                                              0    225b",
		Index:            ".apm-agent-configuration",
		Shard:            "0",
		PriRep:           "r",
		State:            "STARTED",
		Node:             "aragorn",
		UnassignedReason: "",
		Docs:             "0",
		Store:            225,
	},
	{
		TestName: "Relocating",
		Input: "index                                                                    shard prirep state      node                                                        unassigned.reason     docs   store\n" +
			".ds-pods_wisel-gateway---2024.01.23-000025                               0     p      RELOCATING gandalf -> 192.168.168.247 xC9WIDrVRNuiC2CVcu55Ig galadriel                   42981806   5.5gb",
		Index:            ".ds-pods_wisel-gateway---2024.01.23-000025",
		Shard:            "0",
		PriRep:           "p",
		State:            "RELOCATING",
		Node:             "gandalf -> 192.168.168.247 xC9WIDrVRNuiC2CVcu55Ig galadriel",
		UnassignedReason: "",
		Docs:             "42981806",
		Store:            5905580032,
	},
}

func TestShardsParsing(t *testing.T) {
	for _, tc := range TestCases {
		log.Info().Msgf("ShardsParsing Test case: %v", tc.TestName)
		lines := strings.Split(tc.Input, "\n")
		headers := parseHeaderLine(lines[0])
		shard, err := parseLine(lines[1], headers)
		assert.NoError(t, err, "parseLine should not return an error")
		log.Info().Msgf("Shard: %v", shard.ToString())

		assert.Equal(t, tc.Index, shard.Index, "Index should be equal")
		assert.Equal(t, tc.Shard, shard.Shard, "Shard should be equal")
		assert.Equal(t, tc.PriRep, shard.PriRep, "PriRep should be equal")
		assert.Equal(t, tc.State, shard.State, "State should be equal")
		assert.Equal(t, tc.Node, shard.Node, "Node should be equal")
		assert.Equal(t, tc.UnassignedReason, shard.UnassignedReason, "UnassignedReason should be equal")
		assert.Equal(t, tc.Docs, shard.Docs, "Docs should be equal")
		assert.Equal(t, tc.Store, shard.Store, "Store should be equal")
	}
}
