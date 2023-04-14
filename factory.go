package webhookeventreceiver

import (
	"context"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	// The value of "type" key in configuration.
	typeStr = "generic_webhook"
)

type webhookeventreceiverFactory struct {
	receivers *sharedcomponent.SharedComponents
}

// Default configuration for the generic webhook receiver
func (f *webhookeventreceiverFactory) createDefaultConfig() component.Config {
	return &Config{
		HTTPServerSettings: confighttp.HTTPServerSettings{
			Endpoint: "localhost:8080",
		},
		ReadTimeout: 60 * time.Second,
	}
}

func (f *webhookeventreceiverFactory) createLogsReceiver(ctx context.Context, params receiver.CreateSettings, cfg component.Config, consumer consumer.Logs) (r receiver.Logs, err error) {
	rcfg := cfg.(*Config)
	return newwebhookeventreceiverReceiver(rcfg, consumer, params)
}

// NewFactory creates a factory for Generic Webhook Receiver.
func NewFactory() receiver.Factory {
	f := &webhookeventreceiverFactory{
		receivers: sharedcomponent.NewSharedComponents(),
	}
	return receiver.NewFactory(
		typeStr,
		f.createDefaultConfig,
		f.createDefaultConfig,
		receiver.WithLogs(f.createLogsReceiver, component.StabilityLevelAlpha))
}
