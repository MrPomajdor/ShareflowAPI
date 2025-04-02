package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/MrPomajdor/ShareFlowAPI/internal/auth"
	"github.com/MrPomajdor/ShareFlowAPI/internal/config"
	errors "github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	"github.com/MrPomajdor/ShareFlowAPI/internal/healthcheck"
	"github.com/MrPomajdor/ShareFlowAPI/internal/info"
	accesslog "github.com/MrPomajdor/ShareFlowAPI/pkg/accesslog"
	"github.com/MrPomajdor/ShareFlowAPI/pkg/dbcontext"
	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	content "github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/default.yaml", "path to config file")

func main() {
	flag.Parse()
	logger := logrus.New()

	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.WithField("error", err.Error()).Fatal("Failed to load config file")
	}
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Info("Invalid log level")
	}
	logger.SetLevel(level)
	logger.WithField("level", cfg.LogLevel).Info("Set log level")
	db, err := dbx.MustOpen("mysql", cfg.DSN)
	if err != nil {
		logger.WithField("error", err.Error()).Fatal("Failed to connect to the database")

	}

	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)

	defer db.Close()

	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db), cfg),
	}
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.WithFields(logrus.Fields{"verison": Version, "address": address}).Info("Server is running")
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger *logrus.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	rg := router.Group("/v1")

	authHandler := auth.Handler(cfg.JWTSigningKey)

	info.RegisterHandlers(rg.Group(""),
		info.NewService(logger, db),
		authHandler, logger,
	)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, db, logger),
		logger,
	)

	return router
}

func logDBQuery(logger *logrus.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, query string, rows *sql.Rows, err error) {
		if err == nil {
			logger.WithContext(ctx).WithFields(logrus.Fields{"query": query, "duration": t.Milliseconds()}).Debug("Database query succesfull")

		} else {
			logger.WithContext(ctx).WithField("error", err.Error()).Error("Database query error!")
		}
	}
}

func logDBExec(logger *logrus.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, query string, result sql.Result, err error) {
		if err == nil {
			logger.WithContext(ctx).WithFields(logrus.Fields{"query": query, "duration": t.Milliseconds()}).Debug("Database execution succesfull")

		} else {
			logger.WithContext(ctx).WithField("error", err.Error()).Error("Database execution error!")
		}
	}
}
