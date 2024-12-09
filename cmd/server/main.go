package main

import (
	"github.com/xChygyNx/metrical/internal/server"
	"log"
)

func main() {
	err := server.Routing()

	if err != nil {
		log.Fatal(err)
	}
}
