package config

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/shahbaz275817/prismo/pkg/logger"
)

func checkKey(key string) {
	if !viper.IsSet(key) {
		logger.WithContext(context.Background()).Errorf("%s key is not set", key)
		panic(fmt.Sprintf("%s key is not set", key))
	}
}

func panicIfErrorForKey(err error, key string) {
	if err != nil {
		panic(fmt.Sprintf("Could not parse key: %s, Error: %s", key, err))
	}
}

func getStringOrPanic(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func getIntOrPanic(key string) int {
	v, err := strconv.Atoi(getStringOrPanic(key))
	panicIfErrorForKey(err, key)
	return v
}

func getFloatOrPanic(key string) float64 {
	v, err := strconv.Atoi(getStringOrPanic(key))
	panicIfErrorForKey(err, key)
	return float64(v)
}

func getBool(key string) bool {
	v, err := strconv.ParseBool(getStringOrPanic(key))
	if err != nil {
		return false
	}
	return v
}

func getString(key string) string {
	return viper.GetString(key)
}

func splitStringOrPanic(key string, delimiter string) []string {
	checkKey(key)
	return strings.Split(viper.GetString(key), delimiter)
}

func splitString(key string, delimiter string) []string {
	return strings.Split(viper.GetString(key), delimiter)
}

func getStringSliceOrPanic(key string, defaultVal ...[]string) []string {
	if len(defaultVal) > 0 {
		viper.SetDefault(key, defaultVal[0])
	}
	checkKey(key)
	return viper.GetStringSlice(key)
}

func getStringMapOrPanic(key string) map[string]interface{} {
	checkKey(key)
	return viper.GetStringMap(key)
}

func getInt64OrPanic(key string) int64 {
	checkKey(key)
	return viper.GetInt64(key)
}

func getInt(key string) int {
	return viper.GetInt(key)
}

func getInt64(key string) int64 {
	return viper.GetInt64(key)
}

func getDuration(key string) time.Duration {
	checkKey(key)
	return viper.GetDuration(key)
}

func unmarshalIfPresent(s string, h interface{}) error {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}
	return json.Unmarshal([]byte(s), h)
}
