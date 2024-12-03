package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routing() error {
	router := chi.NewRouter()
	router.Post("/update/{mType}/{metric}/{value}", SaveMetricHandle)
	router.Get("/value/{mType}/{metric}", GetMetricHandle)
	router.Get("/", ListMetricHandle)

	config, err := getConfig()
	if err != nil {
		return err
	}

	err = http.ListenAndServe(config.hostAddr.String(), router)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := routing()

	if err != nil {
		panic(err)
	}
}
