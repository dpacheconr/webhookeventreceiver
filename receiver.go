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

// Start spins up the receiver's HTTP server and makes the receiver start its processing.
func (whr *webhookeventreceiver) Start(_ context.Context, host component.Host) error {
	whr.settings.Logger.Info("Webhookreceiver starting")
	if host == nil {
		return errMissingHost
	}

	var err error
	whr.server, err = whr.config.HTTPServerSettings.ToServer(host, whr.settings.TelemetrySettings, whr)
	if err != nil {
		return err
	}

	var listener net.Listener
	listener, err = whr.config.HTTPServerSettings.ToListener()
	if err != nil {
		return err
	}
	whr.shutdownWG.Add(1)
	go func() {
		defer whr.shutdownWG.Done()

		if errHTTP := whr.server.Serve(listener); errHTTP != nil && !errors.Is(errHTTP, http.ErrServerClosed) {
			host.ReportFatalError(errHTTP)
		}
	}()

	return nil
}

// Shutdown tells the receiver that should stop reception,
// giving it a chance to perform any necessary clean-up and
// shutting down its HTTP server.
func (whr *webhookeventreceiver) Shutdown(context.Context) error {
	whr.settings.Logger.Info("Closing Webhookreceiver")
	if whr.server == nil {
		return nil
	}
	err := whr.server.Close()
	whr.shutdownWG.Wait()
	return err
}

// ServeHTTP receives webhookevent events and processes them
func (whr *webhookeventreceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := whr.getBody(r)
	if r.URL.Path != "/webhook" {
		w.WriteHeader(404)
		return
	}
	switch r.Method {
	case "POST":
		// Always reply with 200 OK even if we can't process it
		w.WriteHeader(200)
		// Consume the event received
		whr.consumer.Consume(ctx, body)

	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}

}

func (whr *webhookeventreceiver) getBody(r *http.Request) (body string) {
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
	} else {
		whr.settings.Logger.Info("Body: No Body Supplied\n")
	}
	bodystring := prettyJSON.String()

	return bodystring

}
