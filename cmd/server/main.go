package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/template-go/internal/config"
	"github.com/example/template-go/internal/db"
	"github.com/example/template-go/internal/observability"
	"github.com/example/template-go/internal/todo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	tracerShutdown, err := observability.Setup(ctx, "template-go", "0.1.0")
	if err != nil {
		logger.Warn("otel setup failed", "error", err)
	}
	metrics, err := observability.SetupMetrics()
	if err != nil {
		logger.Warn("metrics setup failed", "error", err)
	} else {
		go func() {
			if err := http.ListenAndServe(cfg.MetricsAddr, metrics.Handler); err != nil {
				logger.Error("metrics server error", "error", err)
			}
		}()
	}

	const (
		maxDBAttempts = 60
		dbRetryDelay  = 5 * time.Second
	)

	var pool *pgxpool.Pool
	err = retry(ctx, maxDBAttempts, dbRetryDelay, func() error {
		var err error
		pool, err = db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			logger.Warn("db open failed", "error", err)
		}
		return err
	})
	if err != nil {
		logger.Error("db open failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	err = retry(ctx, maxDBAttempts, dbRetryDelay, func() error {
		err := db.WithConn(ctx, pool, func(conn *pgx.Conn) error {
			if err := db.Migrate(ctx, conn); err != nil {
				return err
			}
			if cfg.SeedData && cfg.Environment != "gcp" {
				return db.Seed(ctx, conn)
			}
			return nil
		})
		if err != nil {
			logger.Warn("migrate/seed retry", "error", err)
		}
		return err
	})
	if err != nil {
		logger.Error("migrate/seed failed", "error", err)
		os.Exit(1)
	}

	svc, err := todo.NewService(pool, logger)
	if err != nil {
		logger.Error("service init failed", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	path, handler := todo.NewConnectHandler(svc)
	mux.Handle(path, handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("unhealthy"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: otelhttp.NewHandler(mux, "todo-api"),
	}

	go func() {
		logger.Info("todo service running", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdown)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	if metrics != nil {
		_ = metrics.Shutdown(shutdownCtx)
	}
	if tracerShutdown != nil {
		_ = tracerShutdown(shutdownCtx)
	}
}

func retry(ctx context.Context, attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return err
}
