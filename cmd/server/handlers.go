package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseGaugeMetricValue(value string) error {
	_, err := strconv.ParseFloat(value, 64)
	return err
}

func parseCounterMetricValue(value string) error {
	_, err := strconv.Atoi(value)
	return err
}

func parseMetricValue(value string, mType string) (err error) {
	switch mType {
	case "gauge":
		err = parseGaugeMetricValue(value)
	case "counter":
		err = parseCounterMetricValue(value)
	}
	return
}

func MetricHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	metric := chi.URLParam(req, "metric")
	if metric == "" {
		http.Error(res, "Metric param is missed", http.StatusInternalServerError)
		return
	}
	mType := chi.URLParam(req, "mType")
	if found := slices.Contains([]string{"gauge", "counter"}, mType); !found {
		http.Error(res, fmt.Sprintf("Type of metric is %s, must be {gauge, counter}", mType), http.StatusBadRequest)
	}
	valueStr := chi.URLParam(req, "value")
	if valueStr == "" {
		http.Error(res, "Value of metric is missed", http.StatusNotFound)
		return
	}
	err := parseMetricValue(valueStr, mType)
	if err != nil {
		http.Error(res, "Value of metric must be numeric, got "+valueStr, http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]string{
		"status": http.StatusText(http.StatusOK),
		"metric": metric,
		"value":  valueStr,
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(data))
		if err != nil {
			http.Error(res, "Internal error", http.StatusInternalServerError)
		}
	}
}
