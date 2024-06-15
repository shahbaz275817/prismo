package inmemory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bluele/gcache"
	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/logger"
	"github.com/shahbaz275817/prismo/pkg/reporting"
)

type CacheStats struct {
	lookupCount uint64
	hitCount    uint64
	hitRate     float64
	missCount   uint64
}

type InMemCache interface {
	LoadValue(ctx context.Context, key interface{}) (interface{}, error)
	RemoveKey(ctx context.Context, key interface{})
	Stats() CacheStats
}

type LoaderFunc func(interface{}) (interface{}, error)

type CacheConfig struct {
	LoaderFunc                   LoaderFunc
	TTLSeconds                   int64
	Reporter                     *reporting.Reporter
	Name                         string
	MetricsPublishIntervalInSecs int64
	Size                         int
}

type GcCacheInMemCache struct {
	Name       string
	LoaderFunc LoaderFunc
	TTLSeconds int64
	stats      CacheStats
	reporter   *reporting.Reporter
	GcCache    gcache.Cache
}

func NewInMemCache(cfg CacheConfig) (cache InMemCache, err error) {
	if cfg.LoaderFunc == nil {
		return cache, errors.NewUnknownError("Cache config loader function is nil")
	}

	gc := gcache.New(cfg.Size).LRU().LoaderFunc(gcache.LoaderFunc(cfg.LoaderFunc))

	isCacheExpirable := cfg.TTLSeconds > 0

	if isCacheExpirable {
		gc.Expiration(time.Duration(cfg.TTLSeconds) * time.Second)
	}

	cache = GcCacheInMemCache{
		Name:       cfg.Name,
		LoaderFunc: cfg.LoaderFunc,
		TTLSeconds: cfg.TTLSeconds,
		stats:      CacheStats{},
		reporter:   cfg.Reporter,
		GcCache:    gc.Build(),
	}

	metricName := fmt.Sprintf("in_mem_cache_%s", strings.ToLower(cfg.Name))

	defer func() {
		if cfg.Reporter == nil {
			return
		}

		/*  Explanation of Metrics:
			1) HitCount: This is the number of times a requested item was found in the cache.
				A "hit" occurs when the cache successfully returns the data for a given key.

			2) MissCount: In contrast to a hit, a "miss" happens when the data for a requested key is not found in the cache.
				The miss count is the total number of these misses. Misses generally lead to a fallback where data is retrieved from a higher latency source, like a database.

			3) LookupCount: This is the total number of cache lookups, which is the sum of both hit-and-miss counts.
				It represents the total number of times the cache was accessed or queried for data.

			4) HitRate: The hit rate is a ratio that represents the effectiveness of the cache.
				It is calculated as the number of hits divided by the total number of lookups (hits plus misses).
		   		A higher hit rate indicates a more effective cache, as it means a higher proportion of requests were fulfilled
		    	directly from the cache without needing to retrieve data from a slower source.
		*/

		f := func() {
			ent := cfg.Reporter.Report(metricName)
			ent.SetGauge("hit_count", cache.Stats().hitCount)
			ent.SetGauge("miss_count", cache.Stats().missCount)
			ent.SetGauge("lookup_count", cache.Stats().lookupCount)
			ent.SetGauge("hit_rate", cache.Stats().hitRate)

			ent.Publish()
		}

		cfg.Reporter.RegisterPeriodicMetrics(f, time.Duration(cfg.MetricsPublishIntervalInSecs)*time.Second)
	}()

	return cache, nil

}

func (c GcCacheInMemCache) LoadValue(ctx context.Context, key interface{}) (interface{}, error) {

	val, err := c.GcCache.Get(key)

	if err != nil {
		logger.WithContext(ctx).Errorf("Error in loading values from cache for key : %s , error : %s", key, err.Error())
		return nil, err
	}

	return val, nil
}

func (c GcCacheInMemCache) RemoveKey(ctx context.Context, key interface{}) {

	res := c.GcCache.Remove(key)

	if !res {
		logger.WithContext(ctx).Errorf("Error in removing key %s from cache %s", key, c.Name)
		return
	}

	logger.WithContext(ctx).Infof("Removed key %s from cache %s", key, c.Name)
}

func (c GcCacheInMemCache) Stats() CacheStats {
	return CacheStats{
		lookupCount: c.GcCache.LookupCount(),
		hitCount:    c.GcCache.HitCount(),
		hitRate:     c.GcCache.HitRate(),
		missCount:   c.GcCache.MissCount(),
	}
}
