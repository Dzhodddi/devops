package db

import (
	"context"
	"database/sql"
	"time"
	"devops/internal/store"
)

func New(addr string, maxOpenConnections, maxIdleConnections int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), store.QueryTimeOut)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil

}
