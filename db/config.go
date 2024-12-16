package db

import "github.com/0xPolygon/cdk-rpc/config/types"

// Config represents the configuration of the json rpc
type Config struct {
	// Enabled defines if the HTTP server is enabled
	Enabled bool `mapstructure:"Enabled"`

	// Host defines the network adapter that will be used to serve the HTTP requests
	Host string `mapstructure:"Host"`

	// Port defines the port to serve the endpoints via HTTP
	Port int `mapstructure:"Port"`

	// ReadTimeout is the HTTP server read timeout
	// check net/http.server.ReadTimeout and net/http.server.ReadHeaderTimeout
	ReadTimeout types.Duration `mapstructure:"ReadTimeout"`

	// WriteTimeout is the HTTP server write timeout
	// check net/http.server.WriteTimeout
	WriteTimeout types.Duration `mapstructure:"WriteTimeout"`

	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`

	AuthMethodList string `mapstructure:"AuthMethodList"`
}
