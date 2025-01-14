package server

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xChygyNx/metrical/internal/server/types"
)

func isMetricInDB(db *sql.DB, table string, metric_name string) (bool, error) {
	var records *sql.Rows
	var err error
	if table == "gauges" {
		records, err = db.QueryContext(context.Background(),
			"SELECT * FROM gauges WHERE metric_name = $1", metric_name)
	} else if table == "counters" {
		records, err = db.QueryContext(context.Background(),
			"SELECT * FROM counters WHERE metric_name = $1", metric_name)
	}
	if err != nil {
		return false, fmt.Errorf("error in search data in table %s: %w", table, err)
	}
	return records.Next(), nil
}

func writeMetricStorageDB(db *sql.DB, storage *types.MemStorage) error {
	ctx := context.Background()
	for k, v := range storage.GetGauges() {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("error in convert gauge value: %w", err)
		}
		metricExist, err := isMetricInDB(db, "gauges", k)
		if err != nil {
			return fmt.Errorf("error in search gauge metric %s in DB: %w", k, err)
		}
		if metricExist {
			_, err = db.ExecContext(ctx, "UPDATE gauges "+
				"SET value = value  +$1 "+
				"WHERE metric_name = $2;", val, k)
			if err != nil {
				return fmt.Errorf("error in update data in gauges table: %w", err)
			}
		} else {
			_, err = db.ExecContext(ctx, "INSERT INTO gauges(metric_name, value) "+
				"VALUES ($1, $2)", k, val)
			if err != nil {
				return fmt.Errorf("error in update data in gauges table: %w", err)
			}
		}
	}
	for k, v := range storage.GetCounters() {
		fmt.Printf("K = %s, V = %s\n", k, v)
		diff, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error in convert counter value: %w", err)
		}
		metricExist, err := isMetricInDB(db, "counters", k)
		if err != nil {
			return fmt.Errorf("error in search counter metric %s in DB: %w", k, err)
		}
		if metricExist {
			_, err = db.ExecContext(ctx, "UPDATE counters "+
				"SET value = value  +$1 "+
				"WHERE metric_name = $2;", diff, k)
			if err != nil {
				return fmt.Errorf("error in update data in counters table: %w", err)
			}
		} else {
			_, err = db.ExecContext(ctx, "INSERT INTO counters(metric_name, value) "+
				"VALUES ($1, $2)", k, diff)
			if err != nil {
				return fmt.Errorf("error in update data in counters table: %w", err)
			}
		}
	}
	return nil
}
