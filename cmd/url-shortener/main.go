package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/wlcmtunknwndth/REST_API/internal/config"
	"github.com/wlcmtunknwndth/REST_API/internal/http-server/handlers/url/save"
	"github.com/wlcmtunknwndth/REST_API/internal/http-server/middleware/logger"
	"github.com/wlcmtunknwndth/REST_API/internal/lib/logger/sl"
	"github.com/wlcmtunknwndth/REST_API/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	//fmt.Println(cfg)

	var log *slog.Logger = setupLogger(cfg.Env)

	log.Info("starting our service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	//id, err := storage.SaveURL("https://pkg.go.dev/database/sql", "sql_lib")
	//if err != nil {
	//	log.Error("can't save url")
	//}
	//log.Info("id:", strconv.FormatInt(id, 10))
	//url, err := storage.GetURL("sql_lib")
	//if err != nil {
	//	log.Error("can't find this alias ", err)
	//}
	//log.Info(url)
	//
	//if err = storage.DeleteURL("google"); err != nil {
	//	log.Error("can't delete alias: ", err)
	//} / f

	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID)
	//router.Use(middleware.RealIP)
	//router.Use(middleware.Logger) //logger option 2
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	// TODO: run server

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

// might use prettylogger
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
