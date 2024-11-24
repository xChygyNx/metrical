package server

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", gaugeHandle)
	mux.HandleFunc("/update/counter/", counterHandle)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%v\n", MemStorage{})
}
