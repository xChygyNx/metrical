package server

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xChygyNx/metrical/internal/server/types"
)

func writeMetricStorageDB(db *sql.DB, storage *types.MemStorage) error {
	ctx := context.Background()
	for k, v := range storage.GetGauges() {
		_, err := db.ExecContext(ctx, "UPDATE gauges"+
			"SET value = $1"+
			"WHERE metric_name = $2;", v, k)
		if err != nil {
			return fmt.Errorf("error in update data in gauges table: %w", err)
		}
	}
	for k, v := range storage.GetCounters() {
		fmt.Printf("K = %s, V = %s\n", k, v)
		diff, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error in convert counter value: %w", err)
		}
		_, err = db.ExecContext(ctx, "UPDATE counters"+
			"SET value = value + $1"+
			"WHERE metric_name = $2;", diff, k)
		if err != nil {
			return fmt.Errorf("error in update data in counters table: %w", err)
		}
	}
	return nil
}
