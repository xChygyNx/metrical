package types

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type gaugeInsertQuery struct {
	query string
	args  []interface{}
	exec  bool
}

func NewGaugeInsertQuery() gaugeInsertQuery {
	return gaugeInsertQuery{
		exec:  false,
		query: "INSERT INTO gauges(metric_name, value) VALUES ",
		args:  make([]interface{}, 0),
	}
}

func (giq *gaugeInsertQuery) AddRecord(metricName string, metricValue string) {
	giq.exec = true
	numArgs := len(giq.args)
	firstArgOffset := 1
	secondArgOffset := 2
	queryParts := []string{giq.query, fmt.Sprintf("($%d, $%d)",
		numArgs+firstArgOffset, numArgs+secondArgOffset)}
	sep := ", "
	if numArgs == 0 {
		sep = " "
	}
	giq.query = strings.Join(queryParts, sep)
	giq.args = append(giq.args, metricName, metricValue)
}

func (giq *gaugeInsertQuery) ExecInsert(ctx context.Context, tx *sql.Tx) (err error) {
	if giq.exec {
		giq.exec = false
		giq.query += ` ON CONFLICT (metric_name) DO UPDATE SET value = EXCLUDED.value`
		_, err = tx.ExecContext(ctx, giq.query, giq.args...)
		if err != nil {
			return fmt.Errorf("error in insert new records in Gauge Tables of Postgresql: %w", err)
		}
	}
	return nil
}

type counterInsertQuery struct {
	query string
	args  []interface{}
	exec  bool
}

func NewCounterInsertQuery() counterInsertQuery {
	return counterInsertQuery{
		exec:  false,
		query: "INSERT INTO counters(metric_name, value) VALUES ",
		args:  make([]interface{}, 0),
	}
}

func (ciq *counterInsertQuery) AddRecord(metricName string, metricValue string) {
	ciq.exec = true
	numArgs := len(ciq.args)
	firstArgOffset := 1
	secondArgOffset := 2
	queryParts := []string{ciq.query, fmt.Sprintf("($%d, $%d)",
		numArgs+firstArgOffset, numArgs+secondArgOffset)}
	sep := ", "
	if numArgs == 0 {
		sep = " "
	}
	ciq.query = strings.Join(queryParts, sep)
	ciq.args = append(ciq.args, metricName, metricValue)
}

func (ciq *counterInsertQuery) ExecInsert(ctx context.Context, tx *sql.Tx) (err error) {
	if ciq.exec {
		ciq.exec = false
		ciq.query += ` ON CONFLICT (metric_name) DO UPDATE SET value = EXCLUDED.value + counters.value`
		_, err = tx.ExecContext(ctx, ciq.query, ciq.args...)
		if err != nil {
			return fmt.Errorf("error in insert new records in Counter Tables of Postgresql: %w", err)
		}
	}
	return nil
}
