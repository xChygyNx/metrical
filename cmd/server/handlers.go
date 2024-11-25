package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func gaugeHandle(res http.ResponseWriter, req *http.Request) {
	// storeMetrics := MemStorage{}
	// err := req.ParseForm()
	// if err != nil {
	// 	http.Error(res, err.Error(), http.StatusBadRequest)
	// }

	path := req.URL.Path
	if match, err := regexp.MatchString(`^\/update\/gauge\/[a-zA-A]+`, path); err != nil || !match {
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else if !match {
			res.WriteHeader(http.StatusNotFound)
		}
	} else if match, err := regexp.MatchString(`^\/update\/gauge\/[a-zA-Z]+\/(\d+\.\d+|\d+)$`, path); err != nil || !match {
		if err != nil {
			http.NotFound(res, req)
		} else if !match {
			res.WriteHeader(http.StatusBadRequest)
		}
	} else {
		path = strings.Replace(path, "/update/", "", 1)
		path = strings.Replace(path, "/", " ", -1)

		var mType, mName string
		var mValue float64
		_, err := fmt.Sscanf(path, "%s %s %g", &mType, &mName, &mValue)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			// storage := GetMemStorage()
			// _, err = json.Marshal(storage)
			// if err != nil {
			// 	http.Error(res, err.Error(), http.StatusInternalServerError)
			// }
			res.WriteHeader(http.StatusOK)
		}
	}
	res.Header().Set("Content-type", "text/plain")
	res.Write([]byte("Привет, я ничего не умею\n"))
}

func counterHandle(res http.ResponseWriter, req *http.Request) {
	// storeMetrics := MemStorage{}
	// err := req.ParseForm()
	// if err != nil {
	// 	http.Error(res, err.Error(), http.StatusBadRequest)
	// }

	path := req.URL.Path

	if match, err := regexp.MatchString(`^\/update\/counter\/[a-zA-A]+`, path); err != nil || !match {
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else if !match {
			res.WriteHeader(http.StatusNotFound)
		}
	} else if match, err := regexp.MatchString(`^\/update\/counter\/[a-zA-Z]+\/\d+$`, path); err != nil || !match {
		if err != nil {
			http.NotFound(res, req)
		} else if !match {
			res.WriteHeader(http.StatusBadRequest)
		}
	} else {
		path = strings.Replace(path, "/update/", "", 1)
		path = strings.Replace(path, "/", " ", -1)

		var mType, mName string
		var mValue float64
		_, err = fmt.Sscanf(path, "%s %s %g", &mType, &mName, &mValue)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		// } else {
		}
	}
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-type", "text/plain")
	res.Write([]byte("Привет, я ничего не умею"))
}
