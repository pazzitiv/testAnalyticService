package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"testAnalyticService/internal/source/pgsql"
)

//go:generate mockgen -source=analytics.go -destination=./analytics_mock.go -package=internal

type AnalyticsRepository interface {
	Add(ctx context.Context, data AnalyticData) error
}

type AnalyticData struct {
	Headers map[string][]string
	Body    AnalyticBody
}

type AnalyticBody struct {
	Module string
	Type   string
	Event  string
	Name   string
	Data   struct {
		Action string
	}
}

type analyticsRepository struct {
	source pgsql.Source
	logger *zap.Logger
}

func NewAnalyticsRepository(source pgsql.Source, logger *zap.Logger) AnalyticsRepository {
	return &analyticsRepository{
		source: source,
		logger: logger,
	}
}

func (r *analyticsRepository) Add(ctx context.Context, data AnalyticData) error {
	timeStamp := time.Now()

	authHeader, ok := data.Headers["X-Tantum-Authorization"]
	if !ok {
		return fmt.Errorf("no auth header")
	}
	if len(authHeader) != 1 {
		return fmt.Errorf("incorrect auth header")
	}
	userId := authHeader[0]
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("add analytics error: %w", err)
	}

	err = r.source.AddAnalytics(ctx, &pgsql.AnalyticsData{
		UserId:    userId,
		Data:      jsonData,
		TimeStamp: timeStamp,
	})
	if err != nil {
		return fmt.Errorf("add analytics to db error: %w", err)
	}

	return nil
}
