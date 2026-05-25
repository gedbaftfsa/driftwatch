// Package ratelimit implements a token-bucket rate limiter keyed by service
// name, allowing driftwatch to control how frequently drift checks are
// triggered for individual services.
//
// Usage:
//
//	// Allow at most 5 checks per minute per service.
//	limiter := ratelimit.New(12*time.Second, 5)
//
//	if limiter.Allow(serviceName) {
//		// perform drift check
//	} else {
//		log.Printf("rate limited: %s", serviceName)
//	}
//
// Tokens are replenished automatically as time passes. Each key (service
// name) maintains an independent bucket, so throttling one service does
// not affect others.
package ratelimit
