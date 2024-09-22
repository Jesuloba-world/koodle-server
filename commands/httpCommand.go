package command

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"

	userrepo "github.com/Jesuloba-world/koodle-server/repo/user"
	authservice "github.com/Jesuloba-world/koodle-server/services/authService"
	otpservice "github.com/Jesuloba-world/koodle-server/services/otpService"
	senderservice "github.com/Jesuloba-world/koodle-server/services/senderService"
	"github.com/Jesuloba-world/koodle-server/util"

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
	config, err := util.GetConfig()
	if err != nil {
		slog.Error("Error reading config", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))
	router.Use(middleware.Recoverer)

	humaConfig := huma.DefaultConfig("Koodle API", "1.0.0")
	api := humachi.New(router, humaConfig)

	userrepo := userrepo.NewUserRepo(db)

	senderService := senderservice.NewSenderService(config.MsKey, "needle@trial-pr9084z2ev84w63d.mlsender.net", userrepo)

	otpExpirationDuration := time.Minute * 30 // 30 minutes
	otpGenerateTimeLapse := time.Minute * 5   // 5 minutes

	otpService := otpservice.NewOTPService(db, otpExpirationDuration, otpGenerateTimeLapse, senderService)

	authService := authservice.NewAuthService(api, db, otpService, userrepo)
	authService.RegisterRoutes()

	slog.Info("Server starting", "port", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
	return nil
}
