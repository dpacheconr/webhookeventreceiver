package webhookeventreceiver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

var (
	errMissingHost = errors.New("nil host")
)

// The webhookeventreceiver is responsible for using the unmarshaler and the consumer.
type webhookeventconsumer interface {
	// Consume unmarshalls and consumes the records.
	Consume(ctx context.Context, records [][]byte, commonAttributes map[string]string) (int, error)
}

// webhookeventreceiver
type webhookeventreceiver struct {
	settings   receiver.CreateSettings
	config     *Config
	server     *http.Server
	shutdownWG sync.WaitGroup
	consumer   webhookeventconsumer
	logger     *zap.Logger
}

var _ receiver.Logs = (*webhookeventreceiver)(nil)
var _ http.Handler = (*webhookeventreceiver)(nil)

// Start spins up the receiver's HTTP server and makes the receiver start
// its processing.
func (fmr *webhookeventreceiver) Start(_ context.Context, host component.Host) error {
	// fmr.logger.Info("Started webhook listener")
	if host == nil {
		return errMissingHost
	}

	var err error
	fmr.server, err = fmr.config.HTTPServerSettings.ToServer(host, fmr.settings.TelemetrySettings, fmr)
	if err != nil {
		return err
	}

	var listener net.Listener
	listener, err = fmr.config.HTTPServerSettings.ToListener()
	if err != nil {
		return err
	}
	fmr.shutdownWG.Add(1)
	go func() {
		defer fmr.shutdownWG.Done()

		if errHTTP := fmr.server.Serve(listener); errHTTP != nil && !errors.Is(errHTTP, http.ErrServerClosed) {
			host.ReportFatalError(errHTTP)
		}
	}()

	return nil
}

// Shutdown tells the receiver that should stop reception,
// giving it a chance to perform any necessary clean-up and
// shutting down its HTTP server.
func (fmr *webhookeventreceiver) Shutdown(context.Context) error {
	// fmr.logger.Info("Closing webhook listener")
	if fmr.server == nil {
		return nil
	}
	err := fmr.server.Close()
	fmr.shutdownWG.Wait()
	return err
}

// ServeHTTP receives webhookevent requests, and sends them along to the consumer,
func (fmr *webhookeventreceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmr.logger.Info("Serving webhook listener")
	// ctx := r.Context()
	fmt.Fprintf(w, "Serving\n")
}
