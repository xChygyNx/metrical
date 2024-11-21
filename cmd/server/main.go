package main

import (
	// "fmt"
	"net/http"
	"encoding/json"
)

func gaugeHandle(res http.ResponseWriter, req *http.Request) {
	// storeMetrics := &MemStorage{}
	err := req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	reqStr, err := json.Marshal(req.Form)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.Write(reqStr)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", gaugeHandle)
	// mux.HandleFunc("/update/counter/", counterHandle)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%v\n", MemStorage{})
}
