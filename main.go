package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"thaifood-backend/auth"
	"thaifood-backend/handler"
	"thaifood-backend/internal/config"
	"thaifood-backend/mw"
	"thaifood-backend/store"
)

//go:embed schema.sql seed.sql
var sqlFiles embed.FS

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if command != "serve" && command != "migrate" && command != "seed" && command != "admin" {
		return errors.New("usage: thaifood [serve|migrate|seed|admin EMAIL]")
	}
	var c config.Config
	var err error
	if command == "serve" {
		c, err = config.Load()
		if err != nil {
			return err
		}
	} else {
		c.DatabaseURL = os.Getenv("DATABASE_URL")
		if c.DatabaseURL == "" {
			return errors.New("DATABASE_URL is required")
		}
	}
	poolConfig, err := pgxpool.ParseConfig(c.DatabaseURL)
	if err != nil {
		return errors.New("invalid DATABASE_URL configuration")
	}
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 0
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second
	// Avoid named prepared statements for Supabase transaction pooler compatibility.
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec
	poolConfig.ConnConfig.RuntimeParams["statement_timeout"] = "10000"
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		cancel()
		return errors.New("database initialization failed")
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	cancel()
	if err != nil {
		return errors.New("database connection failed; check DATABASE_URL and Supabase project status")
	}
	if command != "serve" {
		return databaseCommand(pool, command)
	}
	h := &handler.Handler{Store: &store.Store{Pool: pool}, Tickets: &auth.Issuer{Secret: []byte(c.JWTSecret), TTL: time.Hour}}
	routes := h.Routes()
	bounded := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		routes.ServeHTTP(w, r.WithContext(ctx))
	})
	server := &http.Server{
		Addr:              ":" + c.Port,
		Handler:           mw.Logging(mw.Recover(mw.CORS(c.AllowedOrigins, bounded))),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20,
	}
	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("ThaiFood Guide API listening on :%s", c.Port)
		serverErr <- server.ListenAndServe()
	}()
	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.New("HTTP server failed")
		}
	case <-stopCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return errors.New("HTTP shutdown timed out")
		}
	}
	return nil
}

func databaseCommand(pool *pgxpool.Pool, command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if command == "admin" {
		if len(os.Args) != 3 {
			return errors.New("usage: thaifood admin EMAIL (register the account first)")
		}
		tag, err := pool.Exec(ctx, `UPDATE members SET role='admin',updated_at=now() WHERE lower(email)=lower($1)`, strings.TrimSpace(os.Args[2]))
		if err != nil {
			return errors.New("admin promotion failed")
		}
		if tag.RowsAffected() != 1 {
			return errors.New("member not found; register the account first")
		}
		log.Print("Registered member promoted to admin")
		return nil
	}
	file := "schema.sql"
	if command == "seed" {
		file = "seed.sql"
	}
	data, err := sqlFiles.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", file, err)
	}
	if _, err = pool.Exec(ctx, string(data)); err != nil {
		return fmt.Errorf("%s failed; check database schema and permissions", command)
	}
	log.Printf("Database %s completed", command)
	return nil
}
