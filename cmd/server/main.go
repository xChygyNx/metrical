package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", BadRequestHandle)
	mux.HandleFunc("/update/gauge/", GaugeHandle)
	mux.HandleFunc("/update/counter/", CounterHandle)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%v\n", MemStorage{})
}
