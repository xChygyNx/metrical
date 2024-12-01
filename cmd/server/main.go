package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func routing() error {
	router := chi.NewRouter()
	router.Post("/update/gauge/", NotFoundHandle)
	router.Post("/update/counter/", NotFoundHandle)
	router.Post("/update/gauge/{metric}/{value}", GaugeHandle)
	router.Post("/update/counter/{metric}/{value}", CounterHandle)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("route does not exist"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	serverAddr := parseFlag()
	serverAddrStr := fmt.Sprintf("%s:%d", serverAddr.Host, serverAddr.Port)

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
