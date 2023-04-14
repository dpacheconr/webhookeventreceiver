package webhookeventreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"

	"go.opentelemetry.io/collector/internal/sharedcomponent"
	// "github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent"
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

var receivers = sharedcomponent.NewSharedComponents[*Config, *webhookeventreceiver]()

// createLog creates a log receiver based on provided config.
func createLogsReceiver(_ context.Context, set receiver.CreateSettings, cfg component.Config, consumer consumer.Logs) (receiver.Logs, error) {
	oCfg := cfg.(*Config)
	r, err := receivers.GetOrAdd(oCfg, func() (*webhookeventreceiver, error) {
		return newwebhookeventreceiverReceiver(oCfg, set)
	})
	if err != nil {
		return nil, err
	}

	if err = r.Unwrap().registerLogsConsumer(consumer); err != nil {
		return nil, err
	}
	return r, nil
}

// NewFactory creates a factory for Generic Webhook Receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha))
}
