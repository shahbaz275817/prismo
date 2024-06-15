package utils

import "strconv"

func FormatFloat32IntoString(floatVal float32) string {
	// Convert float32 to string
	return strconv.FormatFloat(float64(floatVal), 'f', -1, 32)
}
