package utils

import (
	"strconv"
)

func ToFloat32(str string) float32 {
	s, _ := strconv.ParseFloat(str, 32)
	return float32(s)
}

func Float32ToStringValue(str float32) string {
	s := strconv.FormatFloat(float64(str), 'g', 5, 32)
	return s
}

func IntToStringValue(str int) string {
	s := strconv.Itoa(str)
	return s
}

func IntToString(str int) *string {
	s := strconv.Itoa(str)
	return &s
}

func StringToInt(str string) int {
	s, _ := strconv.Atoi(str)
	return s
}
