package webhookeventreceiver

import (
	"context"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
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

	lc := &LogsConsumer{
		consumer: nextConsumer,
	}

	return &webhookeventreceiver{
		settings: set,
		config:   config,
		consumer: lc,
	}, nil
}

func (lc *LogsConsumer) Consume(ctx context.Context, payload string) (int, error) {
	logs := plog.NewLogs()
	rl := logs.ResourceLogs().AppendEmpty()
	resourceAttributes := rl.Resource().Attributes()
	resourceAttributes.PutStr("lc.consumer", "test")

	logRecord := rl.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	logRecord.Body().SetStr(payload)
	lc.consumer.ConsumeLogs(ctx, logs)

	return http.StatusOK, nil
}
