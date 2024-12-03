package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routing() error {
	router := chi.NewRouter()
	router.Post("/update/{mType}/{metric}/{value}", SaveMetricHandle)
	router.Get("/value/{mType}/{metric}", GetMetricHandle)
	router.Get("/", ListMetricHandle)

	serverAddr := parseFlag()
	serverAddrStr := fmt.Sprintf("http://%s:%d", serverAddr.Host, serverAddr.Port)

	err := http.ListenAndServe(serverAddrStr, router)
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
	// fmt.Printf("%v\n", MemStorage{})
}
