package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"testAnalyticService/internal"
	"testAnalyticService/internal/api/http"
	"testAnalyticService/internal/source/pgsql"
)

func main() {
	var cfg Config

	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("can't init logger: %v", err)
		return
	}

	ctx := context.Background()
	wg := sync.WaitGroup{}

	dbConn, err := initDb(ctx, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		logger.Fatal(fmt.Sprintf("init db error: %s", err))
	}

	analysticsSource := pgsql.NewAnalyticsSource(dbConn)
	analysticsRepo := internal.NewAnalyticsRepository(logger, analysticsSource)

	wg.Add(1)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Fatal("http-server panic error", zap.Any("error", e))
			}
			wg.Done()
		}()

		httpServer := http.NewHTTPServer(cfg.Host, cfg.Port, analysticsRepo, logger)
		err := httpServer.Start(ctx)
		if err != nil {
			logger.Fatal("http-server error", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		<-ctx.Done()
		err := dbConn.Close(ctx)
		if err != nil {
			logger.Error("db connect close error", zap.Error(err))
		}
	}()

	wg.Wait()
	logger.Warn("Application is shutdown")
}

func initDb(ctx context.Context, host string, port int, dbName, user, password string) (*pgx.Conn, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbName)
	dbCtx, dbCancel := context.WithCancel(ctx)
	defer dbCancel()

	conn, err := pgx.Connect(dbCtx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connect to database error: %s", err)
	}

	return conn, nil
}

func initLogger(logLevel string) (*zap.Logger, error) {
	opts := zap.NewProductionConfig()
	loglvl := zap.NewAtomicLevel()
	err := loglvl.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, err
	}
	opts.Level = loglvl

	logger, err := opts.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
