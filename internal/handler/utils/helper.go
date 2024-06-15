package utils

import (
	"fmt"
	"time"
)

func GetUnixTimestampIfPresent(t *time.Time) *int64 {
	if t == nil {
		return nil
	}

	return GetIntPointer(t.Unix())
}

func GetUnixMilliTimestampIfPresent(t *time.Time) *int64 {
	if t == nil {
		return nil
	}

	return GetIntPointer(t.UnixMilli())
}

func GetIntPointer(value int64) *int64 {
	return &value
}

func FormatVolumeFloat32(volume *float32) string {
	if volume == nil {
		return ""
	}

	return fmt.Sprintf("%f m^3", *volume)
}
