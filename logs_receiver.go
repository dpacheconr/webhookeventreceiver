package webhookeventreceiver

import (
	"context"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

// The metricsConsumer implements the firehoseConsumer
// to use a metrics consumer and unmarshaler.
type LogsConsumer struct {
	// consumer passes the translated metrics on to the
	// next consumer.
	consumer consumer.Logs
}

var _ webhookeventconsumer = (*LogsConsumer)(nil)

// newMetricsReceiver creates a new instance of the receiver
// with a metricsConsumer.
func newMetricsReceiver(
	config *Config,
	set receiver.CreateSettings,
	nextConsumer consumer.Logs,
) (receiver.Logs, error) {
	if nextConsumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	mc := &LogsConsumer{
		consumer: nextConsumer,
	}

	return &webhookeventreceiver{
		settings: set,
		config:   config,
		consumer: mc,
	}, nil
}

func (mc *LogsConsumer) Consume(ctx context.Context, records [][]byte, commonAttributes map[string]string) (int, error) {
	return http.StatusOK, nil
}
