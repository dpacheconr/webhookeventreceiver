package webhookeventreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	// The value of "type" key in configuration.
	typeStr = "generic_webhook"
)

// Default configuration for the generic webhook receiver
func createDefaultConfig() component.Config {
	return &Config{
		HTTPServerSettings: confighttp.HTTPServerSettings{
			Endpoint: "localhost:8080",
		},
		ReadTimeout: 60 * time.Second,
	}
}

func createLogsReceiver(ctx context.Context, params receiver.CreateSettings, cfg component.Config, consumer consumer.Logs) (r receiver.Logs, err error) {
	rcfg := cfg.(*Config)
	return newwebhookeventreceiverReceiver(rcfg, consumer, params)
}

// NewFactory creates a factory for Generic Webhook Receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha))
}
