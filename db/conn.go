package db

import (
	"database/sql"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func Conn(connString string) (*sql.DB, error) {
	dsn := connString

	db, err := otelsql.Open("postgres", dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{OmitConnResetSession: true}),
		otelsql.WithTracerProvider(otel.GetTracerProvider()),
	)
	if err != nil {
		return nil, err
	}

	_, err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
