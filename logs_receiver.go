package webhookeventreceiver

import (
	"context"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

type LogsConsumer struct {
	consumer consumer.Logs
}

var _ webhookeventconsumer = (*LogsConsumer)(nil)

func newLogsReceiver(
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

func (mc *LogsConsumer) Consume(ctx context.Context) (int, error) {

	return http.StatusOK, nil
}
