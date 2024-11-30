package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func BadRequestHandle(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func GaugeHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	metric := chi.URLParam(req, "metric")
	if metric == "" {
		http.Error(res, "Metric param is missed", http.StatusInternalServerError)
		return
	}
	value := chi.URLParam(req, "value")
	if value == "" {
		http.Error(res, "Value of metric is missed", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]string{"status": http.StatusText(http.StatusOK)})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(data))
	}
}

func CounterHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	metric := chi.URLParam(req, "metric")
	if metric == "" {
		http.Error(res, "Metric param is missed", http.StatusInternalServerError)
		return
	}
	value := chi.URLParam(req, "value")
	if value == "" {
		http.Error(res, "Value of metric is missed", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]string{"status": http.StatusText(http.StatusOK)})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(data))
	}
}
