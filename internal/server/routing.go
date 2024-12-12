package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routing() error {
	storage := GetMemStorage()
	router := chi.NewRouter()
	router.Post("/update/{mType}/{metric}/{value}", SaveMetricHandle(storage))
	router.Get("/value/{mType}/{metric}", GetMetricHandle(storage))
	router.Get("/", ListMetricHandle(storage))

	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("error in GetConfig: %w", err)
	}

	err = http.ListenAndServe(config.HostAddr.String(), router)
	if err != nil {
		return fmt.Errorf("error with launch http server: %w", err)
	}
	return nil
}
