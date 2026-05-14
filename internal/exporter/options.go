package exporter

import "io"

// Option is a functional option for configuring an Exporter.
type Option func(*Config)

// WithFormat sets the output encoding format.
func WithFormat(f Format) Option {
	return func(c *Config) {
		c.Format = f
	}
}

// WithOutput sets the io.Writer the Exporter writes to.
func WithOutput(w io.Writer) Option {
	return func(c *Config) {
		c.Output = w
	}
}

// NewWithOptions creates an Exporter starting from DefaultConfig and applying
// each Option in order.
func NewWithOptions(opts ...Option) *Exporter {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return New(cfg)
}
