package main

import (
	"log"

	"github.com/xChygyNx/metrical/internal/agent"
)

func main() {
	err := agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
