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
	match, err := regexp.MatchString(`^/update/gauge/[a-zA-A]+/(\d+\.\d+|\d+)$`, path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	if !match {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
	}

	path = strings.Replace(path, "/update/", "", 1)
	path = strings.Replace(path, "/", " ", -1)

	var mType, mName string
	var mValue float64
	_, err = fmt.Sscanf(path, "%s %s %g", &mType, &mName, &mValue)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	// storage := GetMemStorage()
	// _, err = json.Marshal(storage)
	// if err != nil {
	// 	http.Error(res, err.Error(), http.StatusInternalServerError)
	// }
	res.Write([]byte("Привет, я ниего не умею"))
}

func counterHandle(res http.ResponseWriter, req *http.Request) {
	// storeMetrics := MemStorage{}
	// err := req.ParseForm()
	// if err != nil {
	// 	http.Error(res, err.Error(), http.StatusBadRequest)
	// }

	path := req.URL.Path
	match, err := regexp.MatchString(`^/update/counter/[a-zA-A]+/(\d+\.\d+|\d+)$`, path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	if !match {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
	}

	path = strings.Replace(path, "/update/", "", 1)
	path = strings.Replace(path, "/", " ", -1)

	var mType, mName string
	var mValue float64
	_, err = fmt.Sscanf(path, "%s %s %g", &mType, &mName, &mValue)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Write([]byte("Привет, я ниего не умею"))
}
