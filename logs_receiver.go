package webhookeventreceiver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
)

type LogsConsumer struct {
	consumer consumer.Logs
}

var _ webhookeventconsumer = (*LogsConsumer)(nil)

/* global variable declaration */
var servicename string

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
	servicename = config.ServiceName

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
	resourceAttributes.PutStr("service.name", servicename)
	logRecord := rl.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()

	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	var result map[string]any
	json.Unmarshal([]byte(payload), &result)
	// logRecord.Body().SetEmptyMap().FromRaw(result)
	logRecord.Body().SetStr(payload)
	lc.consumer.ConsumeLogs(ctx, logs)

	return http.StatusOK, nil
}
