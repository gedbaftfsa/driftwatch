package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/provider"
	"github.com/example/driftwatch/internal/report"
	"github.com/example/driftwatch/internal/snapshot"
)

func main() {
	var (
		configPath  = flag.String("config", "driftwatch.yaml", "path to config file")
		format      = flag.String("format", "text", "output format: text or json")
		snapshotDir = flag.String("snapshot-dir", ".driftwatch/snapshots", "directory for snapshots")
		saveSnap    = flag.Bool("save-snapshot", false, "save drift summary as snapshot")
	)
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rp := provider.NewRuntimeProvider(cfg.RuntimeEndpoint)
	runtime, err := rp.FetchAll(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching runtime state: %v\n", err)
		os.Exit(1)
	}

	results := drift.Detect(cfg.Services, runtime)
	summary := drift.Summarize(results)

	if *saveSnap {
		store, err := snapshot.NewStore(*snapshotDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error initialising snapshot store: %v\n", err)
			os.Exit(1)
		}
		if err := store.Save(summary); err != nil {
			fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
			os.Exit(1)
		}
	}

	fmt := report.NewFormatter(*format)
	output, err := fmt.Render(summary)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering report: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(output)

	if summary.DriftCount > 0 {
		os.Exit(2)
	}
}
