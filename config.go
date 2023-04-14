package webhookeventreceiver

import (
	"time"

	"go.opentelemetry.io/collector/config/confighttp"
)

// Config defines configuration for the Generic Webhook receiver.
type Config struct {
	confighttp.HTTPServerSettings `mapstructure:",squash"`
	// ReadTimeout of the http server
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}
