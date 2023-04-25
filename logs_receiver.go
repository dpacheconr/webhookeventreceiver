package webhookeventreceiver

import (
	"context"
	"encoding/json"
	"fmt"
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
	sl := rl.ScopeLogs().AppendEmpty()
	resourceAttributes := rl.Resource().Attributes()
	resourceAttributes.PutStr("service.name", servicename)
	logRecord := sl.LogRecords().AppendEmpty()
	// logRecord := rl.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()

	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	var result map[string]any
	err := json.Unmarshal([]byte(payload), &result)
	if err != nil {
		// print out if error is not nil
		fmt.Println(err)
	}

	// parse nested json and set is as log attributes
	for key, val := range result {
		if rec, ok := val.(map[string]interface{}); ok {
			for key1, val1 := range rec {
				val1 := fmt.Sprintf("%v", val1)
				resourceAttributes.PutStr(key+"."+key1, val1)
			}
		} else {
			value := fmt.Sprintf("%v", rec)
			fmt.Println(key + "." + value)
		}
	}
	// or set parsed json as log body by updating the below method
	logRecord.Body().SetStr("wenhookeventdata")
	lc.consumer.ConsumeLogs(ctx, logs)

	return http.StatusOK, nil
}
