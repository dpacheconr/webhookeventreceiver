package webhookeventreceiver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
)

var (
	errMissingHost = errors.New("nil host")
)

type webhookeventconsumer interface {
	Consume(ctx context.Context, payload string) (int, error)
}

// webhookeventreceiver
type webhookeventreceiver struct {
	settings   receiver.CreateSettings
	config     *Config
	server     *http.Server
	shutdownWG sync.WaitGroup
	consumer   webhookeventconsumer
}

var _ receiver.Logs = (*webhookeventreceiver)(nil)
var _ http.Handler = (*webhookeventreceiver)(nil)

// Start spins up the receiver's HTTP server and makes the receiver start
// its processing.
func (fmr *webhookeventreceiver) Start(_ context.Context, host component.Host) error {
	fmr.settings.Logger.Info("Webhookreceiver starting")
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
	fmr.settings.Logger.Info("Closing Webhookreceiver")
	if fmr.server == nil {
		return nil
	}
	err := fmr.server.Close()
	fmr.shutdownWG.Wait()
	return err
}

// ServeHTTP receives webhookevent requests, and sends them along to the consumer,
func (fmr *webhookeventreceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Fprintf(w, "served")
	// r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	// dec := json.NewDecoder(r.Body)
	var bodyBytes []byte
	var err error
	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer r.Body.Close()
	}
	var prettyJSON bytes.Buffer
	if len(bodyBytes) > 0 {
		if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
			fmt.Printf("JSON parse error: %v", err)
			return
		}
		fmt.Println(prettyJSON.String())
	} else {
		fmt.Printf("Body: No Body Supplied\n")
	}
	// testbdy := new(bytes.Buffer)
	// testbdy.ReadFrom(r.Body)
	// spew.Dump(testbdy.String())
	fmr.consumer.Consume(ctx, prettyJSON.String())
}

// getBody reads the body from the request as a slice of bytes.
// func (fmr *webhookeventreceiver) getBody(r *http.Request) ([]byte, error) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = r.Body.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return body, nil
// }
