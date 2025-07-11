package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xChygyNx/metrical/internal/agent"
)

// For output tech info need launch app with ldflags
//
//	-X main.buildVersion=<buildVersion> for output "Build version: <buildVersion>" (default N/A)
//	-X main.buildDate=$(date +'<date_format>') for output "Build date: <date_in_format>" (default N/A)
//	-X main.buildCommit=$(git rev-parse HEAD) for output "Build commit: <last_commit_hash>" (default N/A)
var buildVersion string = "N/A" // Server version
var buildDate string = "N/A"    // Date of Build
var buildCommit string = "N/A"  // Hash of the last commit

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	config, err := agent.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	errGroup, ctx := errgroup.WithContext(context.Background())
	errGroup.Go(func() error {
		return agent.RunGRPC(ctx, sigs, config)
	})
	//errGroup.Go(func() error {
	//	return agent.RunHTTP(ctx, sigs, config)
	//})
	if err := errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}
