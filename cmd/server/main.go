package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xChygyNx/metrical/internal/server"
)

func routing() error {
	storage := server.GetMemStorage()
	router := chi.NewRouter()
	router.Post("/update/{mType}/{metric}/{value}", server.SaveMetricHandle(storage))
	router.Get("/value/{mType}/{metric}", server.GetMetricHandle(storage))
	router.Get("/", server.ListMetricHandle(storage))

	config, err := server.GetConfig()
	if err != nil {
		return err
	}

	err = http.ListenAndServe(config.HostAddr.String(), router)
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
