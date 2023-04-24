package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"testAnalyticService/internal"
	"testAnalyticService/internal/api/http"
	"testAnalyticService/internal/source/pgsql"
	"testAnalyticService/internal/worker"
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
	defer func() {
		if e := recover(); e != nil {
			logger.Fatal("panic error", zap.Error(fmt.Errorf("%s", e)))
		}
	}()

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	wg := sync.WaitGroup{}

	dbConn, err := initDb(ctx, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		logger.Fatal(fmt.Sprintf("init db error: %s", err))
	}

	pgSource := pgsql.NewSource(dbConn)
	analysticsRepo := internal.NewAnalyticsRepository(pgSource, logger)
	analysticsWorker := worker.NewWorker(ctx, analysticsRepo, logger)

	wg.Add(1)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Fatal("http-server panic error", zap.Error(fmt.Errorf("%s", e)))
			}
			wg.Done()
		}()

		httpServer := http.NewHTTPServer(cfg.Host, cfg.Port, analysticsWorker, logger)
		err := httpServer.Start(ctx)
		// Отменяем контекст, если HTTP-сервер завершил работу
		cancelCtx()
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
		dbConn.Close()
	}()

	wg.Wait()
	logger.Warn("Application is shutdown")
}

func initDb(ctx context.Context, host string, port int, dbName, user, password string) (*pgxpool.Pool, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbName)
	dbCtx, dbCancel := context.WithCancel(ctx)
	defer dbCancel()

	conn, err := pgxpool.New(dbCtx, dbUrl)
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
