package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func routing() error {
	router := chi.NewRouter()
	router.Post("/", BadRequestHandle)
	router.Post("/update/gauge/{metric}/{value}", GaugeHandle)
	router.Post("/update/counter/{metric}/{value}", CounterHandle)

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
