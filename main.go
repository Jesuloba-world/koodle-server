package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"

	command "github.com/Jesuloba-world/koodle-server/commands"
	"github.com/Jesuloba-world/koodle-server/migrations"
	"github.com/Jesuloba-world/koodle-server/util"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	config, err := util.GetConfig()
	if err != nil {
		slog.Error("Error reading config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	dsn := config.DatabaseUrl
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("Failed to parse dsn", "error", err)
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("Failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	sqldb := stdlib.OpenDBFromPool(pool)
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		slog.Info("Could not connect to the database", "error", err)
		os.Exit(1)
	}

	app := &cli.App{
		Name:  "xoom",
		Usage: "a cli to control the application",

		Commands: []*cli.Command{
			command.NewDBCommand(migrate.NewMigrator(db, migrations.Migrations)),
			command.HttpCommand(db),
		},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("an error occurred", "error", err)
		os.Exit(1)
	}
}
