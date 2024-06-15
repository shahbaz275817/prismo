package wrappers

import (
	goCtx "context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/shahbaz275817/prismo/internal/responder"
	"github.com/shahbaz275817/prismo/pkg/context"
	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/reporting"
)

type DefaultHandlerFunc func(http.ResponseWriter, *http.Request) error

func DefaultWrapper(handler DefaultHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		report := context.GetReporterEntry(ctx)
		if report != nil {
			defer report.Publish()
			report.Attempt()
		}
		defer handlePanic(ctx, w, req, report)

		hsw := NewHTTPStatusWriter(w)
		err := handler(hsw, req)

		if report != nil {
			mr := fmt.Sprintf("http_status_code.%d", hsw.statusCode)
			report.Incr(mr)
			if err != nil {
				report.Failure()
			} else {
				report.Success()
			}
		}
	}
}

func handlePanic(ctx goCtx.Context, w http.ResponseWriter, req *http.Request, report *reporting.ReporterEntry) {
	if r := recover(); r != nil {
		fmt.Println(r)
		fmt.Println("----- Handler panicked, stacktrace: ----- \n" + string(debug.Stack()))
		report.Failure()
		report.Incr("panic")
		responder.WriteError(w, req, errors.NewInternalServerError("internal_server_error", &errors.ErrDetails{
			Message: "Internal Server Error",
			Params:  nil,
		}))

		panic(r)
	}
}
