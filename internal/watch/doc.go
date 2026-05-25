// Package watch provides continuous drift detection by polling the runtime
// provider on a configurable interval and invoking a handler whenever drift
// is detected between declared service specifications and live state.
//
// Basic usage:
//
//	cfg := watch.DefaultConfig(
//		watch.WithInterval(1 * time.Minute),
//		watch.WithServiceFilter("api", "worker"),
//	)
//	w := watch.New(cfg, runtimeProvider, func(results []drift.Result) {
//		fmt.Println("drift detected:", len(results))
//	}, nil)
//
//	if err := w.Run(ctx, declaredSpecs); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package watch
