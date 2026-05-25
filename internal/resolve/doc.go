// Package resolve provides a Resolver that joins declared service definitions
// from configuration with live runtime states fetched via the provider package.
//
// The resulting []ServiceState slice gives callers a unified view of each
// service's desired state alongside what is actually running, making it
// straightforward to pass directly into the drift detector.
//
// Usage:
//
//	rt  := provider.NewRuntimeProvider(endpoint)
//	res := resolve.New(rt)
//	states, err := res.Resolve(ctx, cfg.Services)
//	// states[i].Declared holds the config entry
//	// states[i].Live is nil when the service was not found at runtime
package resolve
