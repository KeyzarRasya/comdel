package mock

import (
	"comdel-backend/internal/config"
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type MockDBLoader struct{
	LoadFunc func() (config.DBConn, error)
}

func (dbl *MockDBLoader) Load() (config.DBConn, error) {
	return dbl.LoadFunc()
}

type MockDBConn struct {
	BeginFunc func(ctx context.Context) (config.DBTx, error);
	QueryRowFunc func(ctx context.Context, sql string, args ...any) config.DBRow
	QueryFunc func(ctx context.Context, sql string, args ...any) (config.DBRows, error)
}

func (mdbc *MockDBConn) QueryRow(ctx context.Context, sql string, args ...any) config.DBRow {
	return mdbc.QueryRowFunc(ctx, sql, args...);
}

func (mdbc *MockDBConn) Begin(ctx context.Context) (config.DBTx, error) {
	return mdbc.BeginFunc(ctx);
}

func (mdbc *MockDBConn) Query(ctx context.Context, sql string, args ...any) (config.DBRows, error) {
	return mdbc.QueryFunc(ctx, sql, args...)
}

type MockDBTx struct {
	CommitFunc func(ctx context.Context) error;
	RollbackFunc func(ctx context.Context) error;
	ExecFunc func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error);
	QueryFunc func(ctx context.Context, sql string, args ...any) (config.DBRows, error);
	QueryRowFunc func(ctx context.Context, sql string, args ...any) config.DBRow
}

func (mdtx *MockDBTx) Commit(ctx context.Context) error {return mdtx.CommitFunc(ctx)}
func (mdtx *MockDBTx) Rollback(ctx context.Context) error {return mdtx.RollbackFunc(ctx)}
func (mdtx *MockDBTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {return mdtx.ExecFunc(ctx, sql, args...)}
func (mdtx *MockDBTx) Query(ctx context.Context, sql string, args ...any) (config.DBRows, error) {return mdtx.QueryFunc(ctx, sql, args...)}
func (mdtx *MockDBTx) QueryRow(ctx context.Context, sql string, args ...any) config.DBRow {return mdtx.QueryRowFunc(ctx, sql, args...)}

type MockDBRows struct {
	NextFunc func() bool;
	ScanFunc func(dest ...any) error;
	ErrFunc func() error;
	CloseFunc func();
}

func (mrows *MockDBRows) Next() bool {return mrows.NextFunc()}
func (mrows *MockDBRows) Scan(dest ...any) error {return mrows.ScanFunc(dest...)}
func (mrows *MockDBRows) Err() error {return mrows.ErrFunc()}
func (mrows *MockDBRows) Close() {mrows.CloseFunc()}

type MockDBRow struct {
	ScanFunc func(dest ...any) error
}

func (mrow *MockDBRow) Scan(dest ...any) error {
	return mrow.ScanFunc(dest...)
}
