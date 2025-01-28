package types

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type gaugeInsertQuery struct {
	exec  bool
	query string
	args  []interface{}
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
	queryParts := []string{giq.query, fmt.Sprintf("($%d, $%d)", numArgs+1, numArgs+2)}
	giq.query = strings.Join(queryParts, ", ")
	giq.args = append(giq.args, metricName)
	giq.args = append(giq.args, metricValue)
}

func (giq *gaugeInsertQuery) ExecInsert(tx *sql.Tx, ctx context.Context) (err error) {
	if giq.exec {
		_, err = tx.ExecContext(ctx, giq.query, giq.args...)
	} else {
		return
	}
	return err
}

type counterInsertQuery struct {
	exec  bool
	query string
	args  []interface{}
}

func NewCounterInsertQuery() gaugeInsertQuery {
	return gaugeInsertQuery{
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
	ciq.args = append(ciq.args, metricName)
	ciq.args = append(ciq.args, metricValue)
}

func (ciq *counterInsertQuery) ExecInsert(tx *sql.Tx, ctx context.Context) (err error) {
	if ciq.exec {
		_, err = tx.ExecContext(ctx, ciq.query, ciq.args...)
	} else {
		return
	}
	return err
}
