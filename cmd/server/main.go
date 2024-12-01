package main

import (
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
	// fmt.Printf("%v\n", MemStorage{})
}
