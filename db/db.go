package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	// DB postgres connection
	conn         *pgxpool.Pool
	ErrNotActive = errors.New("account is not active")
	ErrBalance   = errors.New("insufficient funds")
)

//Close Shutdown connection of postgres
func Close() {
	conn.Close()
}

// Connect creates new connection of postgres with pgx driver
func Connect(connStr string) error {
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return err
	}

	conn, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return err
	}
	return nil
}

//IsNotFound - check error is no rows
func IsNotFound(err error) bool {
	return err == pgx.ErrNoRows
}
