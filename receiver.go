package webhookeventreceiver

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

type webhookeventreceiverReceiver struct {
	config       *Config
	params       receiver.CreateSettings
	nextConsumer consumer.Logs
	server       *http.Server
	tReceiver    *obsreport.Receiver
	logger       *zap.Logger
}

func newwebhookeventreceiverReceiver(config *Config, nextConsumer consumer.Logs, params receiver.CreateSettings) (receiver.Logs, error) {
	if nextConsumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	instance, err := obsreport.NewReceiver(obsreport.ReceiverSettings{LongLivedCtx: false, ReceiverID: params.ID, Transport: "http", ReceiverCreateSettings: params})
	if err != nil {
		return nil, err
	}
	return &webhookeventreceiverReceiver{
		params:       params,
		config:       config,
		nextConsumer: nextConsumer,
		server: &http.Server{
			ReadTimeout: config.ReadTimeout,
			Addr:        config.HTTPServerSettings.Endpoint,
		},
		tReceiver: instance,
	}, nil
}

func (webhookeventreceiverRcvr *webhookeventreceiverReceiver) Start(_ context.Context, host component.Host) error {
	webhookeventreceiverRcvr.logger.Info("webhookeventreceiver start called")
	go func() {
		ddmux := http.NewServeMux()
		ddmux.HandleFunc("/webhook", webhookeventreceiverRcvr.handleLogs)
		webhookeventreceiverRcvr.server.Handler = ddmux
		if err := webhookeventreceiverRcvr.server.ListenAndServe(); err != http.ErrServerClosed {
			host.ReportFatalError(fmt.Errorf("error starting webhook receiver: %w", err))
		}
	}()
	return nil
}

func (webhookeventreceiverRcvr *webhookeventreceiverReceiver) handleLogs(w http.ResponseWriter, req *http.Request) {
	webhookeventreceiverRcvr.logger.Info("webhookeventreceiver handleLogs called")
	obsCtx := webhookeventreceiverRcvr.tReceiver.StartLogsOp(req.Context())
	var err error
	webhookeventreceiverRcvr.tReceiver.EndTracesOp(obsCtx, "NR", 1, err)
}

func (webhookeventreceiverRcvr *webhookeventreceiverReceiver) Shutdown(ctx context.Context) (err error) {
	webhookeventreceiverRcvr.logger.Info("webhookeventreceiver shutdown called")
	return webhookeventreceiverRcvr.server.Shutdown(ctx)
}
