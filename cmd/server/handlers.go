package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func parseGaugeMetricValue(value string) (num float64, err error) {
	num, err = strconv.ParseFloat(value, 64)
	return
}

func parseCounterMetricValue(value string) (num int64, err error) {
	num, err = strconv.ParseInt(value, 10, 64)
	return
}

func saveMetricValue(mType, mName, value string) (err error) {
	storage := GetMemStorage()
	switch mType {
	case "gauge":
		var num float64
		num, err = parseGaugeMetricValue(value)
		if err != nil {
			return
		}
		storage.SetGauge(mName, num)
	case "counter":
		var num int64
		num, err = parseCounterMetricValue(value)
		if err != nil {
			return
		}
		storage.SetConunter(mName, num)
	}
	return
}

func SaveMetricHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	path := req.URL.Path

	match, err := regexp.MatchString(`^\/update\/(gauge|counter)`, path)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(res, "Unknown metric type, must be gauge or counter", http.StatusBadRequest)
		return
	}

	match, err = regexp.MatchString(`^\/update\/(gauge|counter)/[a-zA-Z0-9]+`, path)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(res, "Incorrect metric name, must contains only from alphabetical and numerical symbols", http.StatusNotFound)
		return
	}

	match, err = regexp.MatchString(`^\/update\/(gauge|counter)/[a-zA-Z0-9]+\/(\d+\.\d+|\d+)$`, path)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(res, "Incorrect metric value, must be numerical", http.StatusBadRequest)
		return
	}

	path = strings.Replace(path, "/update/", "", 1)
	path = strings.Replace(path, "/", " ", -1)

	var mType, mName, valueStr string
	_, err = fmt.Sscanf(path, "%s %s %s", &mType, &mName, &valueStr)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}

	err = saveMetricValue(mType, mName, valueStr)
	if err != nil {
		http.Error(res, "Value of metric must be numeric, got "+valueStr, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte("OK"))
	if err != nil {
		http.Error(res, "Internal error", http.StatusInternalServerError)
	}
}

func getMetricValue(mType, mName string) (num interface{}, ok bool) {
	storage := GetMemStorage()
	switch mType {
	case "gauge":
		num, ok = storage.GetGauge(mName)
		return
	case "counter":
		num, ok = storage.GetCounter(mName)
		return
	}
	return
}

func GetMetricHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	path := req.URL.Path

	match, err := regexp.MatchString(`^\/value\/(gauge|counter)`, path)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(res, "Unknown metric type, must be gauge or counter", http.StatusBadRequest)
		return
	}

	match, err = regexp.MatchString(`^\/value\/(gauge|counter)/[a-zA-Z0-9]+`, path)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(res, "Incorrect metric name, must contains only from alphabetical and numerical symbols", http.StatusNotFound)
		return
	}

	path = strings.Replace(path, "/value/", "", 1)
	path = strings.Replace(path, "/", " ", -1)

	var mType, mName string
	_, err = fmt.Sscanf(path, "%s %s", &mType, &mName)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}

	valueInterface, ok := getMetricValue(mType, mName)
	if !ok {
		http.Error(res, "Metric "+mName+" not set", http.StatusNotFound)
		return
	}

	switch valueInterface.(type) {
	case int64:
		_, err = res.Write([]byte(strconv.FormatInt(valueInterface.(int64), 10)))
		if err != nil {
			http.Error(res, "Internal error", http.StatusInternalServerError)
		}
	case float64:
		_, err = res.Write([]byte(strconv.FormatFloat(valueInterface.(float64), 'f', -1, 64)))
		if err != nil {
			http.Error(res, "Internal error", http.StatusInternalServerError)
		}
	}

	res.WriteHeader(http.StatusOK)

}
