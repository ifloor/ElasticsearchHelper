package sharding

import (
	model2 "ElasticsearchHelper/elasticsearch/sharding/model"
	"fmt"
	"github.com/rs/zerolog/log"
)

func parseHeaderLine(header string) []model2.ShardHeaderPositioned {
	var headers []model2.ShardHeaderPositioned
	var currentHeader *TempToken
	for i, char := range header {
		if currentHeader == nil {
			if char != ' ' {
				currentHeader = &TempToken{
					Token:         string(char),
					FirstPosition: i,
				}
			}
		} else {
			if char == ' ' {
				endedHeader := model2.ShardHeaderPositioned{
					Header:        model2.ShardHeader(currentHeader.Token),
					FirstPosition: currentHeader.FirstPosition,
					LastPosition:  i - 1,
				}
				headers = append(headers, endedHeader)
				currentHeader = nil
			} else {
				currentHeader.Token += string(char)
			}
		}
	}

	if currentHeader != nil {
		endedHeader := model2.ShardHeaderPositioned{
			Header:        model2.ShardHeader(currentHeader.Token),
			FirstPosition: currentHeader.FirstPosition,
			LastPosition:  len(header) - 1,
		}
		headers = append(headers, endedHeader)
	}

	return headers
}

func parseLine(line string, headers []model2.ShardHeaderPositioned) (model2.ElasticShard, error) {
	var shard model2.ElasticShard

	consideringHeaderIndex := 0
	consideringHeader := headers[consideringHeaderIndex]
	var buildingToken *TempToken
	lastSetValue := ""
	var lastSetValueWasOnHeader model2.ShardHeader
	for i, char := range line {
		if buildingToken == nil {
			if char != ' ' {
				buildingToken = &TempToken{
					Token:         string(char),
					FirstPosition: i,
				}
			}
		} else {
			if char != ' ' {
				buildingToken.Token += string(char)
			} else {
				lastPosition := i - 1
				allowedToProceed := false
				for !allowedToProceed {
					if buildingToken != nil && (buildingToken.FirstPosition == consideringHeader.FirstPosition || lastPosition == consideringHeader.LastPosition) {
						setValueOnShard(&shard, consideringHeader.Header, buildingToken.Token)
						lastSetValue = buildingToken.Token
						lastSetValueWasOnHeader = consideringHeader.Header
						//
						consideringHeaderIndex++ // jump to the next header
						consideringHeader = headers[consideringHeaderIndex]
						buildingToken = nil
						allowedToProceed = true
					} else {
						// Probably a multi word value or a value that is past the header (header without value)
						if buildingToken.FirstPosition > consideringHeader.LastPosition {
							// It was past the header
							consideringHeaderIndex++ // jump to the next header
							if consideringHeaderIndex < len(headers) {
								consideringHeader = headers[consideringHeaderIndex]
							}
						} else {
							lastSetValue = fmt.Sprintf("%v %v", lastSetValue, buildingToken.Token)
							setValueOnShard(&shard, lastSetValueWasOnHeader, lastSetValue)
							//
							buildingToken = nil
							allowedToProceed = true
						}
					}
				}
			}
		}
	}
	if buildingToken != nil {
		allowedToProceed := false
		lastPosition := len(line) - 1
		for !allowedToProceed {
			if buildingToken != nil && (buildingToken.FirstPosition == consideringHeader.FirstPosition || lastPosition == consideringHeader.LastPosition) {
				setValueOnShard(&shard, consideringHeader.Header, buildingToken.Token)
				lastSetValue = buildingToken.Token
				lastSetValueWasOnHeader = consideringHeader.Header
				//
				consideringHeaderIndex++ // jump to the next header
				if consideringHeaderIndex < len(headers) {
					consideringHeader = headers[consideringHeaderIndex]
				}
				buildingToken = nil
				allowedToProceed = true
			} else {
				// Probably a multi-word value

				// For now, just skip it
				if buildingToken.FirstPosition > consideringHeader.LastPosition {
					// It was past the header
					consideringHeaderIndex++ // jump to the next header
					consideringHeader = headers[consideringHeaderIndex]
				} else {
					lastSetValue = fmt.Sprintf("%v %v", lastSetValue, buildingToken.Token)
					setValueOnShard(&shard, lastSetValueWasOnHeader, lastSetValue)
					//
					buildingToken = nil
					allowedToProceed = true
				}
			}
		}
	}

	return shard, nil
}

func setValueOnShard(shard *model2.ElasticShard, header model2.ShardHeader, value string) {
	switch header {
	case model2.Index:
		shard.Index = value
	case model2.Shard:
		shard.Shard = value
	case model2.PriRep:
		shard.PriRep = value
	case model2.State:
		shard.State = value
	case model2.Node:
		shard.Node = value
	case model2.UnassignedReason:
		shard.UnassignedReason = value
	case model2.Docs:
		shard.Docs = value
	case model2.Store:
		shard.Store = value
	default:
		log.Warn().Msgf("Unknown header: %v", header)
	}
}
