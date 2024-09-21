package command

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"

	authservice "github.com/Jesuloba-world/koodle-server/services/authService"

)

func HttpCommand(db *bun.DB) *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the HTTP server",
		Action: func(ctx *cli.Context) error {
			return startHttpServer(db)
		},
	}
}

func startHttpServer(db *bun.DB) error {
	port := "10001"
	// config, err := util.GetConfig()
	// if err != nil {
	// 	slog.Error("Error reading config", "error", err)
	// 	os.Exit(1)
	// }
	humaConfig := huma.DefaultConfig("Koodle API", "1.0.0")
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))
	router.Use(middleware.Recoverer)

	api := humachi.New(router, humaConfig)

	authService := authservice.NewAuthService(api)
	authService.RegisterRoutes()

	slog.Info("Server starting", "port", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
	return nil
}
