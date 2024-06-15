package config

import (
	"github.com/shahbaz275817/prismo/pkg/reporting"
)

var statsDCfg reporting.StatsDConfig

func StatsD() reporting.StatsDConfig {
	return statsDCfg
}

func initStatsDConfig() {
	statsDCfg = reporting.StatsDConfig{
		Enabled:   getBool("AMPHIBIAN_STATSD_ENABLED"),
		Host:      getStringOrPanic("AMPHIBIAN_STATSD_HOST"),
		Port:      getIntOrPanic("AMPHIBIAN_STATSD_PORT"),
		Namespace: getStringOrPanic("AMPHIBIAN_STATSD_NAMESPACE"),
		Tags:      getStringSliceOrPanic("AMPHIBIAN_STATSD_TAGS"),
	}
}
