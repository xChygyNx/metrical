package main

import (
	"fmt"
	"log"
	_ "net/http/pprof" // требуется для запуска профилировщика
	"os/exec"
	"time"

	"github.com/xChygyNx/metrical/internal/server"
)

func getLastCommit() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		fmt.Println(err)
	}
	commitHash := string(out)
	return commitHash
}

var buildVersion = "1.0.0"
var buildDate = time.Now().Format("02-01-2006")
var buildCommit = getLastCommit()

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	err := server.Routing()

	if err != nil {
		log.Fatal(err)
	}
}
