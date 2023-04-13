package webhookeventreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	// The value of "type" key in configuration.
	typeStr = "generic_webhook"
	// Default endpoints to bind to.
	defaultEndpoint = ":8080"
)

// Default configuration for the generic webhook receiver
func createDefaultConfig() component.Config {
	return &Config{
		HTTPServerSettings: confighttp.HTTPServerSettings{
			Endpoint: defaultEndpoint,
		},
	}
}

// createLogsReceiver creates a logs receiver based on provided config.
func createLogsReceiver(_ context.Context, params receiver.CreateSettings, baseCfg component.Config, consumer consumer.Logs) (receiver.Logs, error) {
	if consumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	logger := params.Logger
	webhookeventreceiverCfg := baseCfg.(*Config)

	logRcvr := &webhookeventreceiverReceiver{
		logger:       logger,
		nextConsumer: consumer,
		config:       webhookeventreceiverCfg,
	}

	return logRcvr, nil
}

// NewFactory creates a factory for Generic Webhook Receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha))
}
