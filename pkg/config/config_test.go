package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/shahbaz275817/prismo/pkg/config"
)

func TestMustGetUint(t *testing.T) {
	key := "MUST_GET_UINT"
	os.Setenv(key, "4")
	v := config.MustGetUint(key)
	assert.Equal(t, uint(4), v)
	os.Unsetenv(key)
}

func TestMustGetInt(t *testing.T) {
	key := "MUST_GET_INT"
	os.Setenv(key, "4")
	v := config.MustGetInt(key)
	assert.Equal(t, 4, v)
	os.Unsetenv(key)
}

func TestMustGetInt64(t *testing.T) {
	key := "MUST_GET_INT64"
	os.Setenv(key, "4")
	v := config.MustGetInt64(key)
	assert.Equal(t, int64(4), v)
	os.Unsetenv(key)
}

func TestMustGetFloat32(t *testing.T) {
	key := "MUST_GET_FLOAT32"
	os.Setenv(key, "4.5")
	v := config.MustGetFloat32(key)
	assert.Equal(t, float32(4.5), v)
	os.Unsetenv(key)
}

func TestMustGetFloat64(t *testing.T) {
	key := "MUST_GET_FLOAT64"
	os.Setenv(key, "4.5")
	v := config.MustGetFloat64(key)
	assert.Equal(t, 4.5, v)
	os.Unsetenv(key)
}

func TestMustGetString(t *testing.T) {
	key := "MUST_GET_STRING"
	os.Setenv(key, "foo")
	v := config.MustGetString(key)
	assert.Equal(t, "foo", v)
	os.Unsetenv(key)
}

func TestMustGetJSON(t *testing.T) {
	type jsonStruct struct {
		Key string
	}

	val := &jsonStruct{}
	key := "GET_JSON"
	os.Setenv(key, "{\"key\":\"value\"}")
	err := config.MustGetJSON(key, val)
	assert.NoError(t, err)
	assert.Equal(t, &jsonStruct{Key: "value"}, val)

	val = &jsonStruct{}
	os.Setenv(key, "invalid_json")
	err = config.MustGetJSON(key, val)
	assert.Error(t, err)
	assert.Equal(t, &jsonStruct{}, val)
	os.Unsetenv(key)
}

func TestGetInt(t *testing.T) {
	key := "GET_INT"
	os.Setenv(key, "4")
	v := config.GetInt(key)
	assert.Equal(t, 4, v)
	os.Unsetenv(key)
	v = config.GetInt(key)
	assert.Equal(t, 0, v)
}

func TestGetInt64(t *testing.T) {
	key := "GET_INT64"
	os.Setenv(key, "4")
	v := config.GetInt64(key)
	assert.Equal(t, int64(4), v)
	os.Unsetenv(key)
	v = config.GetInt64(key)
	assert.Equal(t, int64(0), v)
}

func TestGetUint(t *testing.T) {
	key := "GET_UINT"
	os.Setenv(key, "4")
	v := config.GetUint(key)
	assert.Equal(t, uint(4), v)
	os.Unsetenv(key)
	v = config.GetUint(key)
	assert.Equal(t, uint(0), v)
}

func TestGetString(t *testing.T) {
	key := "GET_STRING"
	os.Setenv(key, "foo")
	v := config.GetString(key)
	assert.Equal(t, "foo", v)
	os.Unsetenv(key)
	v = config.GetString(key)
	assert.Equal(t, "", v)
}

func TestGetStringSlice(t *testing.T) {
	key := "GET_STRING_SLICE"
	os.Setenv(key, "foo,bar")
	v := config.GetStringSlice(key)
	assert.Equal(t, []string{"foo", "bar"}, v)
	os.Unsetenv(key)
	v = config.GetStringSlice(key)
	assert.Equal(t, []string{}, v)
}

func TestGetIntSlice(t *testing.T) {
	key := "GET_INT_SLICE"
	os.Setenv(key, "1,2,3,4,5")
	v := config.GetIntSlice(key)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, v)
	os.Setenv(key, "1,2,A,3")
	v = config.GetIntSlice(key)
	assert.Equal(t, []int{}, v)
	os.Unsetenv(key)
	v = config.GetIntSlice(key)
	assert.Equal(t, []int{}, v)
}

func TestGetInt64Slice(t *testing.T) {
	key := "GET_INT64_SLICE"
	os.Setenv(key, "1,2,3,4,5")
	v := config.GetInt64Slice(key)
	assert.Equal(t, []int64{1, 2, 3, 4, 5}, v)
	os.Setenv(key, "1,2,A,3")
	v = config.GetInt64Slice(key)
	assert.Equal(t, []int64{}, v)
	os.Unsetenv(key)
	v = config.GetInt64Slice(key)
	assert.Equal(t, []int64{}, v)
}

func TestGetFloat32Slice(t *testing.T) {
	key := "GET_FLOAT32_SLICE"
	os.Setenv(key, "1,2,3,4,5")
	v := config.GetFloat32Slice(key)
	assert.Equal(t, []float32{1, 2, 3, 4, 5}, v)
	os.Setenv(key, "1,2,A,3")
	v = config.GetFloat32Slice(key)
	assert.Equal(t, []float32{}, v)
	os.Unsetenv(key)
	v = config.GetFloat32Slice(key)
	assert.Equal(t, []float32{}, v)
}

func TestGetFloat64Slice(t *testing.T) {
	key := "GET_FLOAT64_SLICE"
	os.Setenv(key, "1,2,3,4,5")
	v := config.GetFloat64Slice(key)
	assert.Equal(t, []float64{1, 2, 3, 4, 5}, v)
	os.Setenv(key, "1,2,A,3")
	v = config.GetFloat64Slice(key)
	assert.Equal(t, []float64{}, v)
	os.Unsetenv(key)
	v = config.GetFloat64Slice(key)
	assert.Equal(t, []float64{}, v)
}

func TestGetFeature(t *testing.T) {
	key := "GET_FEATURE"
	os.Setenv(key, "true")
	v := config.GetFeature(key)
	assert.True(t, v)
	os.Setenv(key, "false")
	v = config.GetFeature(key)
	assert.False(t, v)
	os.Unsetenv(key)
	v = config.GetFeature(key)
	assert.False(t, v)
}

func TestGetJSON(t *testing.T) {
	type jsonStruct struct {
		Key string
	}

	val := &jsonStruct{}
	key := "GET_JSON"
	os.Setenv(key, "{\"key\":\"value\"}")
	err := config.GetJSON(key, val)
	assert.NoError(t, err)
	assert.Equal(t, &jsonStruct{Key: "value"}, val)

	val = &jsonStruct{}
	os.Setenv(key, "")
	err = config.GetJSON(key, val)
	assert.NoError(t, err)
	assert.Equal(t, &jsonStruct{}, val)

	val = &jsonStruct{}
	os.Setenv(key, "invalid_json")
	err = config.GetJSON(key, val)
	assert.Error(t, err)
	assert.Equal(t, &jsonStruct{}, val)
	os.Unsetenv(key)
}

func TestMustGetTimeoutInMS(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{name: "It converts the integer into millisecond.", args: args{key: "HODOR_REDIS_CONNECT_TIMEOUT_MS", value: "100"}, want: 100 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.key, tt.args.value)
			defer os.Unsetenv(tt.args.key)
			assert.Equal(t, tt.want, config.MustGetTimeoutInMS(tt.args.key))
		})
	}
}
