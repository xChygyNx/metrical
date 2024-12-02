package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func NotFoundHandle(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}

func GaugeHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	metric := chi.URLParam(req, "metric")
	if metric == "" {
		http.Error(res, "Metric param is missed", http.StatusInternalServerError)
		return
	}
	valueStr := chi.URLParam(req, "value")
	if valueStr == "" {
		http.Error(res, "Value of metric is missed", http.StatusNotFound)
		return
	}
	_, err := strconv.ParseFloat(valueStr, 64)
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

func CounterHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	metric := chi.URLParam(req, "metric")
	if metric == "" {
		http.Error(res, "Metric param is missed", http.StatusInternalServerError)
		return
	}
	valueStr := chi.URLParam(req, "value")
	if valueStr == "" {
		http.Error(res, "Value of metric is missed", http.StatusNotFound)
		return
	}
	_, err := strconv.Atoi(valueStr)
	if err != nil {
		http.Error(res, "Value of metric must be integer, got "+valueStr, http.StatusBadRequest)
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
