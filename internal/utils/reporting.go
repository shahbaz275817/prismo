package utils

import (
	"context"

	ctxwrapper "github.com/shahbaz275817/prismo/pkg/context"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

func ReportMetrics(ctx context.Context, metric string) {
	reporter := ctxwrapper.GetReporter(ctx)
	if reporter == nil {
		logger.WithContext(ctx).Errorf("unable to report %s metric, reporter instance is nil", metric)
		return
	}
	report := reporter.Report("custom_metric")
	report.Incr(metric)
	report.Publish()
}
