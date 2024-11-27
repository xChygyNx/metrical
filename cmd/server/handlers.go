package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func badRequestHandle(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func gaugeHandle(res http.ResponseWriter, req *http.Request) {
		
	res.Header().Set("Content-type", "text/plain")
	path := req.URL.Path
	match, err := regexp.MatchString(`^\/update\/gauge\/[a-zA-A]+`, path)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	} 
	if !match {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	match, err = regexp.MatchString(`^\/update\/gauge\/[a-zA-Z]+\/(\d+\.\d+|\d+)$`, path)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	} 
	if !match {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	path = strings.Replace(path, "/update/", "", 1)
	path = strings.Replace(path, "/", " ", -1)

	var mType, mName string
	var mValue float64
	_, err = fmt.Sscanf(path, "%s %s %g", &mType, &mName, &mValue)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
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

func counterHandle(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-type", "text/plain")
	path := req.URL.Path

	match, err := regexp.MatchString(`^\/update\/counter\/[a-zA-A]+`, path)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	} 
	if !match {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	
	match, err = regexp.MatchString(`^\/update\/counter\/[a-zA-Z]+\/\d+$`, path)
	if err != nil {
		http.NotFound(res, req)
		return
	} else if !match {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	path = strings.Replace(path, "/update/", "", 1)
	path = strings.Replace(path, "/", " ", -1)

	var mType, mName string
	var mValue int64
	_, err = fmt.Sscanf(path, "%s %s %d", &mType, &mName, &mValue)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
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
