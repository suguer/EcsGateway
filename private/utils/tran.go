package utils

import (
	"strconv"
)

func ToFloat32(str string) float32 {
	s, _ := strconv.ParseFloat(str, 32)
	return float32(s)
}
