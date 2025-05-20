package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/xChygyNx/metrical/internal/server"
)

func main() {
	err := server.Routing()

	if err != nil {
		log.Fatal(err)
	}
}
