package reporting

import (
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"time"
)

type MetricReporter interface {
	Incr(key string)
	Gauge(key string, value interface{})
	getClient() *statsd.Client
}

type Reporter struct {
	MetricReporter MetricReporter
}

type ReporterEntry struct {
	mr           MetricReporter
	parent       string
	incrMetrics  []string
	gaugeMetrics map[string]interface{}
}

func (reporter *Reporter) Report(parent string) *ReporterEntry {
	if reporter == nil {
		return nil
	}

	return &ReporterEntry{
		mr:           reporter.MetricReporter,
		parent:       parent,
		incrMetrics:  []string{},
		gaugeMetrics: map[string]interface{}{},
	}
}

func (reporter *Reporter) GetClient() *statsd.Client {
	return reporter.MetricReporter.getClient()
}

func (entry *ReporterEntry) Success() {
	entry.incrementCounter("success")
}

func (entry *ReporterEntry) Attempt() {
	entry.incrementCounter("attempt")
}

func (entry *ReporterEntry) Timeout() {
	entry.incrementCounter("timeout")
}

func (entry *ReporterEntry) Incr(metric string) {
	entry.incrementCounter(metric)
}

// Publish publishes all the metrics collected by the entry to statsd asynchronously.
// This method should be called by deferring.
func (entry *ReporterEntry) Publish() {
	if entry == nil || entry.mr == nil {
		return
	}

	go func() {
		for _, metric := range entry.incrMetrics {
			entry.mr.Incr(metric)
		}
	}()

	go func() {
		for key, value := range entry.gaugeMetrics {
			entry.mr.Gauge(key, value)
		}
	}()
}

func (entry *ReporterEntry) Failure() {
	entry.incrementCounter("failure")
}

func (entry *ReporterEntry) incrementCounter(counter string) {
	if entry == nil || entry.mr == nil {
		return
	}

	metricKey := fmt.Sprintf("counters.%s.%s.count", entry.parent, counter)
	entry.incrMetrics = append(entry.incrMetrics, metricKey)
}

func (entry *ReporterEntry) SetGauge(gauge string, value interface{}) {
	if entry == nil || entry.mr == nil {
		return
	}
	if entry.gaugeMetrics == nil {
		entry.gaugeMetrics = map[string]interface{}{}
	}

	metricKey := fmt.Sprintf("counters.%s.%s.gauge", entry.parent, gauge)
	entry.gaugeMetrics[metricKey] = value
}

func (reporter *Reporter) RegisterPeriodicMetrics(f func(), period time.Duration) {
	if reporter == nil {
		return
	}

	go func() {
		for {
			time.Sleep(period)
			f()
		}
	}()
}
