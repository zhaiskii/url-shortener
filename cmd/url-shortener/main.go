package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shorteener/internal/config"
	"url-shorteener/internal/http-server/handlers/redirect"
	"url-shorteener/internal/http-server/handlers/url/save"
	nwLogger "url-shorteener/internal/http-server/middleware/logger"
	"url-shorteener/internal/lib/logger/sl"
	"url-shorteener/internal/storage/sqlite"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	envLocal = "local"
	envDev	 = "dev"
	envProd	 = "prod"
)

func main() {
	cfg := config.MustLoad()

	// fmt.Println(cfg)

	log := setupLogger(cfg.Env) 

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(nwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
	})
	router.Get("/{alias}", redirect.New(log, storage))

	server := &http.Server{
		Addr:			cfg.Address,
		Handler:		router,
		ReadTimeout:	cfg.HTTPServer.Timeout,
		WriteTimeout:	cfg.HTTPServer.Timeout,
		IdleTimeout:	cfg.HTTPServer.IdleTimout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to run server")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}