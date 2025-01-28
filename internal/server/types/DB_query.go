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
	secondArgOffset := 1
	queryParts := []string{giq.query, fmt.Sprintf("($%d, $%d)",
		numArgs+firstArgOffset, numArgs+secondArgOffset)}
	giq.query = strings.Join(queryParts, ", ")
	giq.args = append(giq.args, metricName, metricValue)
}

func (giq *gaugeInsertQuery) ExecInsert(ctx context.Context, tx *sql.Tx) (err error) {
	_, err = tx.ExecContext(ctx, giq.query, giq.args...)
	if err != nil {
		return fmt.Errorf("error in insert new records in Gauge Tables of Postgresql: %w", err)
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
		query: "INSERT INTO counter(metric_name, value) VALUES ",
		args:  make([]interface{}, 0),
	}
}

func (ciq *counterInsertQuery) AddRecord(metricName string, metricValue string) {
	ciq.exec = true
	numArgs := len(ciq.args)
	queryParts := []string{ciq.query, fmt.Sprintf("($%d, $%d)", numArgs+1, numArgs+2)}
	ciq.query = strings.Join(queryParts, ", ")
	ciq.args = append(ciq.args, metricName, metricValue)
}

func (ciq *counterInsertQuery) ExecInsert(ctx context.Context, tx *sql.Tx) (err error) {
	_, err = tx.ExecContext(ctx, ciq.query, ciq.args...)
	if err != nil {
		return fmt.Errorf("error in insert new records in Counter Tables of Postgresql: %w", err)
	}
	return nil
}
