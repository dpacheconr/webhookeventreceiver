package webhookeventreceiver

import (
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config defines configuration for the Generic Webhook receiver.
type Config struct {
	confighttp.HTTPServerSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
}
