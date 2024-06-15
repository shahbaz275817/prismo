package inmemory

import (
	"context"
	"testing"
	"time"

	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/reporting"
	"github.com/stretchr/testify/suite"
)

type InMemCacheClientTest struct {
	suite.Suite
	ctx context.Context
}

func (suite *InMemCacheClientTest) SetupSuite() {
	suite.ctx = context.Background()
}

func TestInMemCacheClientTest(t *testing.T) {
	suite.Run(t, new(InMemCacheClientTest))
}

// Test case 1: Test that the function `NewInMemCache()` returns an error when the `LoaderFunc` field in the `CacheConfig` struct is nil.
func (suite *InMemCacheClientTest) TestErrorWhenLoaderFuncIsNil() {

	cfg := CacheConfig{
		LoaderFunc: nil,
		TTLSeconds: 0,
		Reporter:   nil,
		Name:       "test",
		Size:       100,
	}
	_, err := NewInMemCache(cfg)

	suite.Assert().NotNil(err)

}

// Test case 2: Test that the `NewInMemCache()` function returns a `GcCacheInMemCache` struct with the correct values for fields `Name`, `LoaderFunc`, `TTLSeconds`, `stats`, and `reporter` when passed a valid `CacheConfig` struct.
func (suite *InMemCacheClientTest) TestConfigValues() {

	cfg := CacheConfig{
		LoaderFunc: func(key interface{}) (interface{}, error) {
			return key, nil
		},
		TTLSeconds:                   10,
		Reporter:                     &reporting.Reporter{MetricReporter: nil},
		Name:                         "test",
		Size:                         100,
		MetricsPublishIntervalInSecs: 1,
	}
	cache, _ := NewInMemCache(cfg)
	gcCache := cache.(GcCacheInMemCache)
	suite.Assert().Equal("test", gcCache.Name)
	suite.Assert().NotNil(gcCache.GcCache)
	suite.Assert().Equal(int64(10), gcCache.TTLSeconds)
	suite.Assert().NotNil(gcCache.stats)
	suite.Assert().NotNil(gcCache.reporter)
}

// Test case 3: Test that the `LoadValue()` function returns the correct value for a given key when the key is found in the cache.
func (suite *InMemCacheClientTest) TestCorrectLoadValueWhenKeyIsPresent() {

	cfg := CacheConfig{
		LoaderFunc: func(key interface{}) (interface{}, error) {
			return key, nil
		},
		TTLSeconds: 10,
		Reporter:   nil,
		Name:       "test",
		Size:       100,
	}
	cache, _ := NewInMemCache(cfg)
	val, _ := cache.LoadValue(context.Background(), "test_key")
	suite.Assert().Equal("test_key", val)

}

// Test case 4: Test that the `LoadValue()` function returns an error when the key is not found in the cache.
func (suite *InMemCacheClientTest) TestErrorWhenKeyIsAbsent() {

	cfg := CacheConfig{
		LoaderFunc: func(key interface{}) (interface{}, error) {
			return nil, errors.NotFoundError{}
		},
		TTLSeconds: 10,
		Reporter:   nil,
		Name:       "test",
		Size:       100,
	}
	cache, _ := NewInMemCache(cfg)
	_, err := cache.LoadValue(context.Background(), "nonexistent_key")
	suite.Assert().NotNil(err)
}

// Test case 5: Test that the Stats() function returns the correct CacheStats struct for the cache.
func (suite *InMemCacheClientTest) TestCachesMetricsValues() {

	cfg := CacheConfig{
		LoaderFunc: func(key interface{}) (interface{}, error) {
			return key, nil
		},
		TTLSeconds:                   1,
		Reporter:                     &reporting.Reporter{MetricReporter: nil},
		MetricsPublishIntervalInSecs: 1,
		Name:                         "test",
		Size:                         100,
	}
	cache, _ := NewInMemCache(cfg)
	// First call - Cache Miss
	_, err1 := cache.LoadValue(context.Background(), "test_key")
	suite.Assert().Nil(err1)
	// Second call - Cache Hit
	_, err2 := cache.LoadValue(context.Background(), "test_key")
	suite.Assert().Nil(err2)

	stats := cache.Stats()

	time.Sleep(1 * time.Second)
	suite.Assert().Equal(uint64(2), stats.lookupCount)
	suite.Assert().Equal(uint64(1), stats.hitCount)
	suite.Assert().Equal(float64(0.5), stats.hitRate)
	suite.Assert().Equal(uint64(1), stats.missCount)
}

// Test case 6: Test that the `RemoveKey()` function doesn't return error whether Key Removal is success or not
func (suite *InMemCacheClientTest) TestNoErrorWhenRemoveKeyIsCalled() {

	loadedFromCache := 0

	cfg := CacheConfig{
		LoaderFunc: func(key interface{}) (interface{}, error) {
			value, ok := key.(string)
			if ok {
				loadedFromCache++
				return value, nil
			}
			return nil, nil
		},
		TTLSeconds: 10,
		Reporter:   nil,
		Name:       "test",
		Size:       100,
	}
	cache, _ := NewInMemCache(cfg)

	// Loading a value from cache
	newKey := "nonexistent_key"
	value, err := cache.LoadValue(suite.ctx, newKey)
	suite.Assert().Nil(err)
	suite.Equal(1, loadedFromCache)
	suite.Equal(newKey, value)

	// Assert value is removed and cacheLoader is being called twice
	cache.RemoveKey(context.Background(), newKey)
	value, err = cache.LoadValue(suite.ctx, newKey)
	suite.Assert().Nil(err)
	suite.Equal(newKey, value)
	suite.Equal(2, loadedFromCache)
}
