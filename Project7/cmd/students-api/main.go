package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sahilsnghai/golang/Project7/internal/config"
	"github.com/sahilsnghai/golang/Project7/internal/handlers/students"
	"github.com/sahilsnghai/golang/Project7/internal/storage/sqlite"
)

func main() {
	// Load Config
	cfg := config.Mustload()

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// SetUp router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", students.New(storage))
	router.HandleFunc("GET /api/students/{id}", students.GetbyId(storage))
	router.HandleFunc("GET /api/students", students.GetLists(storage))

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("fail to start server")
		}

	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cnl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cnl()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
