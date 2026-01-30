package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/clearstatus/backend/internal/config"
	"github.com/clearstatus/backend/internal/handler"
	"github.com/clearstatus/backend/internal/logger"
	"github.com/clearstatus/backend/internal/middleware"
	"github.com/clearstatus/backend/internal/store"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log/slog"
)

func main() {
	// Load .env from current directory (optional; ignored if file missing)
	_ = godotenv.Load()

	log := logger.New()
	log.Info("starting server")

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Error("config validation failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("config loaded", slog.String("port", cfg.Port))

	db, err := sql.Open(cfg.DBDriver, cfg.DBDSN)
	if err != nil {
		log.Error("database open failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Error("database ping failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	log.Info("database connected", slog.String("db", cfg.ConnectionInfoForLog()))

	dbx := sqlx.NewDb(db, cfg.DBDriver)
	if err := runMigrations(db); err != nil {
		log.Error("migrations failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("migrations completed")

	userStore := store.NewUserStore(dbx)
	orgStore := store.NewOrganizationStore(dbx)
	teamStore := store.NewTeamStore(dbx)
	serviceStore := store.NewServiceStore(dbx)
	incidentStore := store.NewIncidentStore(dbx)

	authHandler := handler.NewAuthHandler(cfg, userStore)
	orgHandler := handler.NewOrganizationHandler(orgStore)
	teamHandler := handler.NewTeamHandler(orgStore, teamStore)
	serviceHandler := handler.NewServiceHandler(orgStore, teamStore, serviceStore)
	incidentHandler := handler.NewIncidentHandler(orgStore, incidentStore)
	publicHandler := handler.NewPublicHandler(orgStore, serviceStore, incidentStore)

	r := chi.NewRouter()
	r.Use(middleware.AllowPreflight) // handle OPTIONS and CORS headers (e.g. frontend :3000 → API :8080)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(log))

	r.Post("/api/auth/google", authHandler.Google)

	r.Route("/api/orgs", func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/", orgHandler.List)
		r.Post("/", orgHandler.Create)
		r.Route("/{orgID}", func(r chi.Router) {
			r.Get("/", orgHandler.Get)
			r.Put("/", orgHandler.Update)
			r.Delete("/", orgHandler.Delete)
			r.Get("/members", orgHandler.ListMembers)
			r.Post("/members", orgHandler.AddMember)
			r.Delete("/members/{userID}", orgHandler.RemoveMember)
			r.Route("/teams", func(r chi.Router) {
				r.Get("/", teamHandler.List)
				r.Post("/", teamHandler.Create)
				r.Get("/{teamID}", teamHandler.Get)
				r.Put("/{teamID}", teamHandler.Update)
				r.Delete("/{teamID}", teamHandler.Delete)
			})
			r.Route("/services", func(r chi.Router) {
				r.Get("/", serviceHandler.List)
				r.Post("/", serviceHandler.Create)
				r.Get("/{serviceID}", serviceHandler.Get)
				r.Put("/{serviceID}", serviceHandler.Update)
				r.Delete("/{serviceID}", serviceHandler.Delete)
			})
			r.Route("/incidents", func(r chi.Router) {
				r.Get("/", incidentHandler.List)
				r.Post("/", incidentHandler.Create)
				r.Get("/{incidentID}", incidentHandler.Get)
				r.Put("/{incidentID}", incidentHandler.Update)
				r.Delete("/{incidentID}", incidentHandler.Delete)
				r.Post("/{incidentID}/updates", incidentHandler.AddUpdate)
			})
		})
	})

	r.Get("/api/public/orgs/{slug}/status", publicHandler.Status)

	addr := ":" + cfg.Port
	log.Info("listening", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Error("server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func runMigrations(db *sql.DB) error {
	migrationsDir := "migrations"
	if d := os.Getenv("MIGRATIONS_DIR"); d != "" {
		migrationsDir = d
	}
	for _, name := range []string{"001_init.up.sql", "002_services_team_id.up.sql"} {
		if name == "002_services_team_id.up.sql" {
			var skip int
			if err := db.QueryRow("SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'services' AND COLUMN_NAME = 'team_id' LIMIT 1").Scan(&skip); err == nil {
				continue // team_id already exists, skip
			}
		}
		path := filepath.Join(migrationsDir, name)
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err = db.Exec(string(b)); err != nil {
			return fmt.Errorf("run migration %s: %w", name, err)
		}
	}
	return nil
}
