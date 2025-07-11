package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	_ "net/http/pprof" // требуется для запуска профилировщика
	"os"
	"os/signal"
	"syscall"

	"github.com/xChygyNx/metrical/internal/server"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	config, err := server.GetConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("error in GetConfig: %w", err).Error())
	}

	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		err := server.RunHTTPServer(ctx, sigs, config)
		if err != nil {
			return fmt.Errorf("error in work HTTP server: %w", err)
		}
		return err
	})
	errGroup.Go(func() error {
		err := server.RunGRPCServer(ctx, sigs, config)
		if err != nil {
			return fmt.Errorf("error in work gRPC server: %w", err)
		}
		return err
	})
	if err := errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}
