package pgsql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=analytics.go -destination=./analytics_mock.go -package=pgsql

type AnalyticsData struct {
	UserId    string    `db:"user_id"`
	Data      []byte    `db:"data"`
	TimeStamp time.Time `db:"time"`
}

type AnalyticsSource interface {
	AddAnalytics(ctx context.Context, data *AnalyticsData) error
}

func NewAnalyticsSource(db *pgx.Conn) AnalyticsSource {
	return &analyticsSource{
		db: db,
	}
}

type analyticsSource struct {
	db *pgx.Conn
}

func (s *analyticsSource) AddAnalytics(ctx context.Context, data *AnalyticsData) error {
	rows := [][]interface{}{{data.UserId, data.Data, data.TimeStamp}}
	_, err := s.db.CopyFrom(
		context.Background(),
		pgx.Identifier{"user_actions"},
		[]string{"user_id", "data", "time"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return nil
}
