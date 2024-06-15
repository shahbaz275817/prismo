package handler

import (
	"net/http"

	"github.com/shahbaz275817/prismo/internal/appcontext/server"
	"github.com/shahbaz275817/prismo/internal/handler/account"
	"github.com/shahbaz275817/prismo/internal/handler/transaction"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	logconst "github.com/shahbaz275817/prismo/constants/fields"
	contextWrapper "github.com/shahbaz275817/prismo/pkg/context"
	"github.com/shahbaz275817/prismo/pkg/logger"

	"github.com/shahbaz275817/prismo/internal/middleware"
)

const (
	frameOptionsKey       = "X-Frame-Options"
	contentTypeKey        = "Content-Type"
	transferEncodingKey   = "Transfer-Encoding"
	xssProtectionKey      = "X-XSS-Protection"
	contentTypeOptionsKey = "X-Content-Type-Options"
	cacheControlKey       = "Cache-Control"
)

func NewRouter(deps server.Dependencies) http.Handler {
	router := mux.NewRouter()

	router.UseEncodedPath()

	router.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	router.Handle("/ping", PingHandler()).Methods(http.MethodGet)

	appRouter := router.PathPrefix("/prismo").Subrouter()
	appRouter.Use(middleware.WithHTTPAuth)

	// Account Handlers
	appRouter.Handle("/v1/accounts", account.CreateAccountHandler(deps.AccountService, deps.AtomicLock)).Methods(http.MethodPost)
	appRouter.Handle("/v1/accounts/{account_id}", account.GetAccountHandler(deps.AccountService)).Methods(http.MethodGet)

	// Transaction Handlers
	appRouter.Handle("/v1/transactions", transaction.CreateTransactionHandler(deps.TransactionService, deps.OperationTypesService, deps.AccountService)).Methods(http.MethodPost)

	newRouter := withAccessLog(withDefaultResponseHeaders(router))
	return http.HandlerFunc(newRouter.ServeHTTP)
}

func withDefaultResponseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Add(frameOptionsKey, "SAMEORIGIN")
		wr.Header().Add(contentTypeKey, "application/json; charset=utf-8")
		wr.Header().Add(transferEncodingKey, "chunked")
		wr.Header().Add(xssProtectionKey, "1; mode=block")
		wr.Header().Add(contentTypeOptionsKey, "nosniff")
		wr.Header().Add(cacheControlKey, "max-age=0, private, must-revalidate")

		// Set X-Request-ID
		rid := req.Header.Get(logconst.RequestIDKey)
		if rid == "" {
			rid = uuid.New().String()
			req.Header.Set(logconst.RequestIDKey, rid)
		}

		wr.Header().Add(logconst.RequestIDKey, rid)
		req = req.WithContext(contextWrapper.WithValue(req.Context(), logconst.RequestIDKey, rid))
		req = req.WithContext(contextWrapper.WithValue(req.Context(), logconst.RequestURI, req.RequestURI))
		req = req.WithContext(contextWrapper.WithValue(req.Context(), logconst.RequestMethod, req.Method))

		next.ServeHTTP(wr, req)
	})
}

func withAccessLog(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(logger.GetAccessLogFile(), next)
}
