package context

import (
	"context"

	"github.com/shahbaz275817/prismo/pkg/reporting"
)

const (
	APIEndpoint    = contextKey("endpoint")
	metricReporter = contextKey("metricReporter")
	reporterKey    = contextKey("reporter")
	languageKey    = contextKey("language")
)

func GetAPIEndpoint(ctx context.Context) string {
	a, ok := ctx.Value(APIEndpoint).(string)
	if !ok {
		return ""
	}

	return a
}

type contextKey string

func (ckt contextKey) String() string {
	return string(ckt)
}

func GetReporter(ctx context.Context) *reporting.Reporter {
	r, ok := ctx.Value(metricReporter).(*reporting.Reporter)
	if !ok {
		return nil
	}

	return r
}

func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, languageKey, language)
}

func Language(ctx context.Context) string {
	h, ok := ctx.Value(languageKey).(string)
	if !ok {
		h = "id"
	}
	return h
}

func GetReporterEntry(ctx context.Context) *reporting.ReporterEntry {
	rep, ok := ctx.Value(reporterKey).(*reporting.ReporterEntry)
	if !ok {
		return nil
	}
	return rep
}

func WithValue(ctx context.Context, key, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}
