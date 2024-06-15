package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shahbaz275817/prismo/pkg/logger"
)

// ShutdownTimeout for graceful shutdown.
const ShutdownTimeout = 1 * time.Second

// ReadWriteTimeout for specifies timeout for http server read and write.
const ReadWriteTimeout = 10 * time.Second

// New initializes the server with all the routes
func New(router http.Handler) *Server {
	return &Server{
		APIServer: &http.Server{
			Handler:      router,
			ReadTimeout:  ReadWriteTimeout,
			WriteTimeout: ReadWriteTimeout,
		},
	}
}

// Server is a wrapper around http.Server and provides
// Serve method with graceful-shutdown enabled
type Server struct {
	APIServer *http.Server
}

// Serve starts the server and blocks until any termination
// signals and performs graceful shutdown.
func (s *Server) Serve(addr string) {
	s.APIServer.Addr = addr
	go s.listenServer()
	s.waitForShutdown()
}

func (s *Server) listenServer() {
	logger.WithContext(context.Background()).Infof("starting API server in %s", s.APIServer.Addr)
	if err := s.APIServer.ListenAndServe(); err != http.ErrServerClosed {
		logger.WithContext(context.Background()).Fatalf(err.Error())
	}
}

func (s *Server) waitForShutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	<-sig
	logger.WithContext(context.Background()).Infof("API server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	s.APIServer.Shutdown(ctx)
	logger.WithContext(ctx).Infof("API server shutdown complete")
}
