package main

import (
	"github.com/xChygyNx/metrical/internal/agent"
	"log"
)

func main() {
	err := agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
