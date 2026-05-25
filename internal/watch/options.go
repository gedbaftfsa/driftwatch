package watch

import "time"

// Option is a functional option for configuring a Watcher.
type Option func(*Config)

// WithInterval sets the polling interval.
func WithInterval(d time.Duration) Option {
	return func(c *Config) {
		if d > 0 {
			c.Interval = d
		}
	}
}

// WithServiceFilter restricts watching to specific service names.
func WithServiceFilter(names ...string) Option {
	return func(c *Config) {
		c.ServiceNames = append(c.ServiceNames, names...)
	}
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(opts ...Option) Config {
	cfg := Config{
		Interval: 30 * time.Second,
	}
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
