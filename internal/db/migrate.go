package db

import (
	"context"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrate applies embedded SQL migrations in lexical order.
func Migrate(ctx context.Context, conn *pgx.Conn) error {
	entries, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(entries)

	for _, file := range entries {
		content, err := migrationsFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}
		if _, err := conn.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("exec %s: %w", file, err)
		}
	}
	return nil
}

// Seed inserts basic tasks for local testing.
func Seed(ctx context.Context, conn *pgx.Conn) error {
	const q = `
        INSERT INTO todo_tasks (title, description, completed)
        VALUES ($1, $2, $3)
        ON CONFLICT DO NOTHING;
    `
	samples := []struct {
		title       string
		description string
	}{
		{"Plan template", "Map out requirements and interfaces"},
		{"Wire observability", "Ensure otel export works"},
	}
	for _, s := range samples {
		if _, err := conn.Exec(ctx, q, s.title, s.description, false); err != nil {
			return fmt.Errorf("seed todo: %w", err)
		}
	}
	return nil
}

// WithConn obtains a dedicated connection for migrations.
func WithConn(ctx context.Context, pool *pgxpool.Pool, fn func(*pgx.Conn) error) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire: %w", err)
	}
	defer conn.Release()
	return fn(conn.Conn())
}
