package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xChygyNx/metrical/internal/server/types"
)

func isMetricInDB(tx *sql.Tx, table string, metricName string) (bool, error) {
	var records *sql.Rows
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if table == "gauges" {
		records, err = tx.QueryContext(ctx, "SELECT * FROM gauges WHERE metric_name = $1", metricName)
	} else if table == "counters" {
		records, err = tx.QueryContext(ctx, "SELECT * FROM counters WHERE metric_name = $1", metricName)
	}

	if err != nil {
		return false, fmt.Errorf("error in search data in table %s: %w", table, err)
	}
	err = records.Err()
	if err != nil {
		return false, fmt.Errorf("error in search data in table %s: %w", table, err)
	}
	return records.Next(), nil
}

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
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("error in convert gauge value: %w", err)
		}
		metricExist, err := isMetricInDB(tx, "gauges", k)
		if err != nil {
			return fmt.Errorf("error in search gauge metric %s in DB: %w", k, err)
		}
		if metricExist {
			_, err = tx.ExecContext(ctx, "UPDATE gauges "+
				"SET value = $1 "+
				"WHERE metric_name = $2;", val, k)
			if err != nil {
				return fmt.Errorf("error in update data in gauges table: %w", err)
			}
		} else {
			giq.AddRecord(k, v)
		}
	}
	err = giq.ExecInsert(ctx, tx)
	if err != nil {
		return fmt.Errorf("error in execution new record in Gauge metric table in PostgreSQL: %w", err)
	}
	ciq := types.NewCounterInsertQuery()
	for k, v := range storage.GetCounters() {
		log.Printf("K = %s, V = %s\n", k, v)
		diff, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error in convert counter value: %w", err)
		}
		metricExist, err := isMetricInDB(tx, "counters", k)
		if err != nil {
			return fmt.Errorf("error in search counter metric %s in DB: %w", k, err)
		}
		if metricExist {
			_, err = tx.ExecContext(ctx, "UPDATE counters "+
				"SET value = value  +$1 "+
				"WHERE metric_name = $2;", diff, k)
			if err != nil {
				return fmt.Errorf("error in update data in counters table: %w", err)
			}
		} else {
			ciq.AddRecord(k, v)
		}
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
