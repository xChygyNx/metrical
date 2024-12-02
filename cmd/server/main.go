package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routing() error {
	router := chi.NewRouter()
	router.Post("/update/{mType}/{metric}/{value}", MetricHandle)

	err := http.ListenAndServe(":8080", router)
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
