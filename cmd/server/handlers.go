package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func gaugeHandle(res http.ResponseWriter, req *http.Request) {
	// storeMetrics := MemStorage{}
	// err := req.ParseForm()
	// if err != nil {
	// 	http.Error(res, err.Error(), http.StatusBadRequest)
	// }

	path := req.URL.Path
	parts := strings.Split(strings.TrimLeft(path, "/"), "/")
	
	if len(parts) != 4 {
		panic("Not define all parts")
	}

	parsedValue, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metricValue := gauge(parsedValue)
	// storeMetrics.gauges[parts[2]] = metricValue

	res.Write([]byte(fmt.Sprintf("%v", float64(metricValue))))
}

func counterHandle(res http.ResponseWriter, req *http.Request) {
	// storeMetrics := MemStorage{}
	// err := req.ParseForm()
	// if err != nil {
	// 	http.Error(res, err.Error(), http.StatusBadRequest)
	// }

	path := req.URL.Path
	parts := strings.Split(strings.TrimLeft(path, "/"), "/")
	
	if len(parts) != 4 {
		panic("Not define all parts")
	}

	parsedValue, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metricValue := counter(parsedValue)
	// storeMetrics.gauges[parts[2]] = metricValue

	res.Write([]byte(fmt.Sprintf("%v", float64(metricValue))))
}
