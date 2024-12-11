package main

import (
	"log"

	"github.com/xChygyNx/metrical/internal/server"
)

func main() {
	err := server.Routing()

	if err != nil {
		log.Fatal(err)
	}
}
