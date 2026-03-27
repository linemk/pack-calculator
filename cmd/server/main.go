package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/dig"

	"github.com/linemk/pack-calculator/internal/calculator"
	"github.com/linemk/pack-calculator/internal/handler"
	mw "github.com/linemk/pack-calculator/internal/middleware"
	"github.com/linemk/pack-calculator/internal/store"
)

//go:embed web
var webFS embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	c := dig.New()
	must(c.Provide(store.New))
	must(c.Provide(calculator.New))
	must(c.Provide(handler.New))

	if err := c.Invoke(func(h *handler.Handler) error { return run(h, logger) }); err != nil {
		logger.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func run(h *handler.Handler, logger *slog.Logger) error {
	r := mux.NewRouter()
	r.Use(mw.RequestID)
	r.Use(mw.NewLogger(logger))
	r.Use(mw.Recoverer)
	r.Use(mw.CORS)

	h.RegisterRoutes(r)

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		return fmt.Errorf("embed fs: %w", err)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.FS(sub)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server started", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return fmt.Errorf("shutdown: %w", srv.Shutdown(shutdownCtx))
}
