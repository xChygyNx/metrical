package main

import (
	"fmt"
	"github.com/xChygyNx/metrical/internal/agent"
	"log"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	err := agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
