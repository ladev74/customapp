package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"customapp/internal/api/handlers"
	cconfig "customapp/internal/config"
	llogger "customapp/internal/logger"
)

const (
	pathToConfigFile = "./config/config.env"
	shutdownTime    = 30 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	rtp := flag.Float64("rtp", 0, "")
	flag.Parse()
	if *rtp == 0 {
		log.Fatal("rtp flag must be specified")
	} else if *rtp > 1.0 || *rtp < 0.0 {
		log.Fatalf("invalid rtp flag value, ...: %f", *rtp)
	}

	config, err := cconfig.New(pathToConfigFile)
	if err != nil {
		log.Fatalf("cannot initialize config: %v", err)
	}

	logger, err := llogger.New(&config.Logger)
	if err != nil {
		log.Fatalf("cannot initialize logger: %v", err)
	}
	defer logger.Sync()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(llogger.MiddlewareLogger(logger, &config.Logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/get", handlers.Get(*rtp, logger))

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.HttpServer.Host, config.HttpServer.Port),
		Handler: router,
	}

	go func() {
		logger.Info("starting http server", zap.String("addr", srv.Addr))
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("cannot start http server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	logger.Info("received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTime)
	defer shutdownCancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("cannot shutdown http server", zap.Error(err))
		return
	}

	logger.Info("stopping http server", zap.String("addr", srv.Addr))

	logger.Info("application shutdown completed successfully")

}
