package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xChygyNx/metrical/internal/server/types"
)

func writeMetricStorageDB(db *sql.DB, storage *types.MemStorage) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error in create transaction for DB: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	giq := types.NewGaugeInsertQuery()
	for k, v := range storage.GetGauges() {
		giq.AddRecord(k, v)
	}
	err = giq.ExecInsert(ctx, tx)
	if err != nil {
		return fmt.Errorf("error in execution new record in Gauge metric table in PostgreSQL: %w", err)
	}
	ciq := types.NewCounterInsertQuery()
	for k, v := range storage.GetCounters() {
		ciq.AddRecord(k, v)
	}
	err = ciq.ExecInsert(ctx, tx)
	if err != nil {
		return fmt.Errorf("error in execution new record in Counter metric table in PostgreSQL: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error in commit transaction to DB: %w", err)
	}
	return nil
}

func retryDBWrite(db *sql.DB, storage *types.MemStorage, retryCount int) (err error) {
	var pgErr *pgconn.PgError
	delays := make([]time.Duration, 0, retryCount)
	delays = append(delays, 0*time.Second)
	for i := 1; i < retryCount; i++ {
		delays = append(delays, time.Duration(2*i-1)*time.Second)
	}

	for i := 0; i < retryCount; i++ {
		time.Sleep(delays[i])
		err = writeMetricStorageDB(db, storage)
		if err == nil || !errors.As(err, &pgErr) {
			break
		}
	}
	return
}
