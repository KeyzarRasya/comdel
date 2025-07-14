package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

// DBLoader = object yg bikin koneksi
type DBLoader interface {
	Load() (DBConn, error)
}

// DBConn = koneksi pool yang bisa mulai transaksi
type DBConn interface {
	Begin(ctx context.Context) (DBTx, error)
	QueryRow(ctx context.Context, sql string, args ...any) DBRow
	Query(ctx context.Context, sql string, args ...any) (DBRows, error)
}

// DBTx = transaksi yang bisa commit, rollback, exec query, dan return rows
type DBTx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (DBRows, error)
	QueryRow(ctx context.Context, sql string, args ...any) DBRow
}

type DBRow interface {
	Scan(...any) error;
}

// DBRows = hasil query untuk scan data
type DBRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}
