package reporting

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/shahbaz275817/prismo/pkg/errors"
)

type StatsDConfig struct {
	Host      string
	Port      int
	Enabled   bool
	Namespace string
	Tags      []string
}

func (cfg StatsDConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

func NewStatsD(cfg StatsDConfig) (*StatsD, error) {
	if !cfg.Enabled {
		return &StatsD{}, nil
	}

	client, err := statsd.New(statsd.Address(cfg.Address()), statsd.Prefix(cfg.Namespace))
	if err != nil {
		panic(errors.Wrap(err, "failed to initiate StatsD"))
	}
	return &StatsD{client: client}, nil
}

func NewHystrixStatsDCollector(statsDConfig StatsDConfig) {
	if !statsDConfig.Enabled {
		return
	}

	address := fmt.Sprintf("%s:%d", statsDConfig.Host, statsDConfig.Port)
	c, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: address,
		Prefix:     statsDConfig.Namespace + ".hystrix",
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to initiate StatsD for Hystrix"))
	}

	metricCollector.Registry.Register(c.NewStatsdCollector)
}

type StatsD struct {
	client *statsd.Client
}

func (reporter *StatsD) Incr(key string) {
	if reporter.client != nil {
		reporter.client.Increment(key)
	}
}

func (reporter *StatsD) Gauge(key string, value interface{}) {
	if reporter.client != nil {
		reporter.client.Gauge(key, value)
	}
}

func (reporter *StatsD) getClient() *statsd.Client {
	return reporter.client
}
