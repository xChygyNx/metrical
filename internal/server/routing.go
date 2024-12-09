package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routing() error {
	storage := GetMemStorage()
	router := chi.NewRouter()
	router.Post("/update/{mType}/{metric}/{value}", SaveMetricHandle(storage))
	router.Get("/value/{mType}/{metric}", GetMetricHandle(storage))
	router.Get("/", ListMetricHandle(storage))

	config, err := GetConfig()
	if err != nil {
		return err
	}

	err = http.ListenAndServe(config.HostAddr.String(), router)
	if err != nil {
		return err
	}
	return nil
}
