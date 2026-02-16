package main

import (
	"boucherie-api/configs"
	"boucherie-api/internal/handler"
	mw "boucherie-api/internal/middleware"
	"boucherie-api/internal/repository"
	"boucherie-api/internal/service"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

func main() {
	// ── Logger ──────────────────────────────────────────
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	// ── Config ──────────────────────────────────────────
	cfg := configs.Load()
	log.Info().Int("port", cfg.Port).Str("db", cfg.DBPath).Msg("configuration loaded")

	// ── Database ────────────────────────────────────────
	db, err := sql.Open("sqlite", cfg.DBPath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}

	// Run schema migrations
	if err := runMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	log.Info().Msg("database schema ready")

	// Run seed data (idempotent via INSERT OR IGNORE)
	if err := runSeed(db); err != nil {
		log.Warn().Err(err).Msg("seed data warning (may already exist)")
	} else {
		log.Info().Msg("seed data loaded")
	}

	// ── Repositories ────────────────────────────────────
	clientRepo := repository.NewClientRepo(db)
	productRepo := repository.NewProductRepo(db)
	saleRepo := repository.NewSaleRepo(db)
	creditRepo := repository.NewCreditRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	// ── Services ────────────────────────────────────────
	clientSvc := service.NewClientService(clientRepo)
	productSvc := service.NewProductService(productRepo)
	saleSvc := service.NewSaleService(saleRepo, productRepo, clientRepo, creditRepo)
	creditSvc := service.NewCreditService(creditRepo, clientRepo)
	orderSvc := service.NewOrderService(orderRepo, clientRepo, productRepo)

	// ── Handlers ────────────────────────────────────────
	clientH := handler.NewClientHandler(clientSvc)
	productH := handler.NewProductHandler(productSvc)
	saleH := handler.NewSaleHandler(saleSvc)
	creditH := handler.NewCreditHandler(creditSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	dashboardH := handler.NewDashboardHandler(db)

	// ── Router ──────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware
	r.Use(mw.CORS())
	r.Use(mw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		handler.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Handle("/dashboard", dashboardH)
		r.Mount("/clients", clientH.Routes())
		r.Mount("/products", productH.Routes())
		r.Mount("/sales", saleH.Routes())
		r.Mount("/credits", creditH.Routes())
		r.Mount("/orders", orderH.Routes())
	})

	// ── Start ───────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Info().Str("addr", addr).Msg("🥩 Boucherie API démarrée")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

// runMigrations executes the schema.sql file to set up tables.
func runMigrations(db *sql.DB) error {
	schema, err := os.ReadFile("migrations/schema.sql")
	if err != nil {
		return fmt.Errorf("reading schema: %w", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("executing schema: %w", err)
	}
	return nil
}

// runSeed loads seed data from migrations/seed.sql (idempotent).
func runSeed(db *sql.DB) error {
	seed, err := os.ReadFile("migrations/seed.sql")
	if err != nil {
		return fmt.Errorf("reading seed: %w", err)
	}
	_, err = db.Exec(string(seed))
	if err != nil {
		return fmt.Errorf("executing seed: %w", err)
	}
	return nil
}
