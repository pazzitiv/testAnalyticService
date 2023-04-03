package pgsql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=source.go -destination=./source_mock.go -package=pgsql

type AnalyticsData struct {
	UserId    string    `db:"user_id"`
	Data      []byte    `db:"data"`
	TimeStamp time.Time `db:"time"`
}

type Source interface {
	AddAnalytics(ctx context.Context, data *AnalyticsData) error
}

func NewSource(db *pgxpool.Pool) Source {
	return &source{
		db: db,
	}
}

type source struct {
	db *pgxpool.Pool
}

func (s *source) AddAnalytics(ctx context.Context, data *AnalyticsData) error {
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
