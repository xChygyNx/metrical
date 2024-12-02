package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func MetricHandle(res http.ResponseWriter, req *http.Request) {

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

	match, err = regexp.MatchString(`^\/update\/(gauge|counter)/[a-zA-Z]+`, path)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(res, "Incorrect metric name, must contains only from alphabetical symbols", http.StatusNotFound)
		return
	}

	match, err = regexp.MatchString(`^\/update\/(gauge|counter)/[a-zA-Z]+\/(\d+\.\d+|\d+)$`, path)
	if err != nil {
		http.Error(res, "Internal Server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(res, "Incorrect metric value, must be numerical", http.StatusNotFound)
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

	data, err := json.Marshal(map[string]string{
		"status":      http.StatusText(http.StatusOK),
		"metricType":  mType,
		"metricName":  mName,
		"metricValue": valueStr,
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
