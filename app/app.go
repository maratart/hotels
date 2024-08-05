package app

import (
	"context"
	"errors"
	"hotels/app/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg Config
	log *logger.Logger
	mux http.Handler
}

func NewApp(cfg Config) App {
	log := logger.NewLogger()
	handlers := NewHandlers(log)
	return App{
		cfg,
		log,
		newMux(handlers),
	}
}

func (a App) Run() {
	a.log.Info("running webserver...")
	a.log.Info("server listening on localhost:%s", a.cfg.port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:         ":" + a.cfg.port,
		Handler:      a.mux,
		ReadTimeout:  a.cfg.readTimeout,
		WriteTimeout: a.cfg.writeTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	// Start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Errorf("server failed: %s", err)
			os.Exit(1)
		} else {
			a.log.Info("server closed")
		}
	}()

	<-ctx.Done()
	a.log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.log.Errorf("server shutdown failed: %s", err)
		os.Exit(1)
	} else {
		a.log.Info("server shutdown gracefully")
	}
}
