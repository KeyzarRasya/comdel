package config

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBLoaderImpl struct{
	
}

func (dl *DBLoaderImpl) Load() (DBConn, error) {
	var databaseUri string;

	if os.Getenv("DEV_ENV") == "dev" {
		databaseUri = os.Getenv("LOCAL_DATABASE_URI")
	} else {
		databaseUri = os.Getenv("DATABASE_URI")
	}

	conn, err := pgxpool.New(context.Background(), databaseUri)

	if err != nil {
		log.Info(os.Getenv("DATABASE_URI"))
		
		log.Error("Failed to Load Database");
		log.Error(err.Error())
		return nil, err;
	}

	log.Info("Creating Database Connection");
	return &DBConnImpl{Pool: conn}, nil;
}

type DBConnImpl struct{
	Pool *pgxpool.Pool

}

func (dbc *DBConnImpl) Begin(context context.Context) (DBTx, error) {
	tx, err := dbc.Pool.Begin(context);

	if err != nil {
		return nil, err;
	}

	return &DBTxImpl{Tx: tx}, nil
}

func (dbc *DBConnImpl) QueryRow(ctx context.Context, sql string, args ...any) DBRow {
	return dbc.Pool.QueryRow(ctx, sql, args...)
}

func (dbc *DBConnImpl) Query(ctx context.Context, sql string, args ...any) (DBRows, error) {
	return dbc.Pool.Query(ctx, sql, args...)
}

type DBTxImpl struct {
	Tx pgx.Tx
}

func (dbt *DBTxImpl) Commit(ctx context.Context) error {
	return dbt.Tx.Commit(ctx)
}

func (dbt *DBTxImpl) Rollback(ctx context.Context) error {return dbt.Tx.Rollback(ctx)}
func (dbt *DBTxImpl) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {return dbt.Tx.Exec(ctx, sql, args...)}
func (dbt *DBTxImpl) Query(ctx context.Context, sql string, args ...any) (DBRows, error) {return dbt.Tx.Query(ctx, sql, args)}
func (dbt *DBTxImpl) QueryRow(ctx context.Context, sql string, args ...any) DBRow {return dbt.Tx.QueryRow(ctx, sql, args...)}

type DBRowImpl struct {
	Row pgx.Row
}

func (dbr *DBRowImpl) Scan(dest ...any) error {
	return dbr.Row.Scan(dest...)
}

type DBRowsImpl struct {
	Rows pgx.Rows
}

	// Next() bool
	// Scan(dest ...any) error
	// Err() error
	// Close()

func (dbrs *DBRowsImpl) Next() bool {return dbrs.Rows.Next()}
func (dbrs *DBRowsImpl) Scan(dest ...any) error {return dbrs.Rows.Scan(dest...)}
func (dbrs *DBRowsImpl) Err() error {return dbrs.Rows.Err()}
func (dbrs *DBRowsImpl) Close() {dbrs.Close()}


