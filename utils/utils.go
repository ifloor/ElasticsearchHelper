package utils

import (
	"strconv"
	"strings"
)

type Dimension struct {
	Dimension string
	factor    int64
}

var B = Dimension{Dimension: "b", factor: 1}
var KB = Dimension{Dimension: "kb", factor: 1024}
var MB = Dimension{Dimension: "mb", factor: 1024 * 1024}
var GB = Dimension{Dimension: "gb", factor: 1024 * 1024 * 1024}
var TB = Dimension{Dimension: "tb", factor: 1024 * 1024 * 1024 * 1024}
var PB = Dimension{Dimension: "pb", factor: 1024 * 1024 * 1024 * 1024 * 1024}

func StringBytesToIntBytes(stringBytes string) int64 {

	var dimension Dimension
	if strings.HasSuffix(stringBytes, KB.Dimension) {
		dimension = KB
	} else if strings.HasSuffix(stringBytes, MB.Dimension) {
		dimension = MB
	} else if strings.HasSuffix(stringBytes, GB.Dimension) {
		dimension = GB
	} else if strings.HasSuffix(stringBytes, TB.Dimension) {
		dimension = TB
	} else if strings.HasSuffix(stringBytes, PB.Dimension) {
		dimension = PB
	} else if strings.HasSuffix(stringBytes, B.Dimension) {
		dimension = B
	} else {
		panic("Impossible to find dimension for string bytes: " + stringBytes + ", better to not continue...")
	}
	stringFloat := strings.ReplaceAll(stringBytes, dimension.Dimension, "")
	valueInDimension, err := strconv.ParseFloat(stringFloat, 64)
	if err != nil {
		panic("Impossible to convert string to float: " + stringFloat + ", better to not continue...")
	}

	return int64(valueInDimension * float64(dimension.factor))
}
