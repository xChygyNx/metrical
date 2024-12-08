package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xChygyNx/metrical/internal/server"
)

func routing() error {
	router := chi.NewRouter()
	router.Post("/update/{mType}/{metric}/{value}", server.SaveMetricHandle)
	router.Get("/value/{mType}/{metric}", server.GetMetricHandle)
	router.Get("/", server.ListMetricHandle)

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
