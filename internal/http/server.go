package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"social-network-api/internal/redis"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
	cache  *redis.Client
}

func New(logger *zap.SugaredLogger, db *pgxpool.Pool, cache *redis.Client) *Server {
	return &Server{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func (s Server) Run() {
	router := s.setHTTPRouter()

	srv := &http.Server{
		Addr:         "0.0.0.0:5000",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)

		// Listen for SIGNINT and SIGTERM signals
		// and write them in quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		s.logger.Infoln("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)

		if err != nil {
			shutdownError <- err
		}

		s.logger.Infoln("Completing background tasks...")

		shutdownError <- nil
	}()

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal(err)
	}

	err = <-shutdownError
	if err != nil {
		s.logger.Fatal(err)
	}

	s.logger.Infoln("Stopped server")
}
