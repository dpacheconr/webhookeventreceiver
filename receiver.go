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

// type Record struct {
// 	Content string `json:"content"`
// }

type Bird struct {
	Species     string
	Description string
}

var _ receiver.Logs = (*webhookeventreceiver)(nil)
var _ http.Handler = (*webhookeventreceiver)(nil)

// Start spins up the receiver's HTTP server and makes the receiver start its processing.
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

// ServeHTTP receives webhookevent events and processes them
func (fmr *webhookeventreceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := fmr.getBody(r)
	// reqBody, _ := ioutil.ReadAll(r.Body)
	// fmt.Fprintf(w, "%+v", string(reqBody))
	fmr.consumer.Consume(ctx, body)
}

func (fmr *webhookeventreceiver) getBody(r *http.Request) (body string) {
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
		fmr.settings.Logger.Info(prettyJSON.String())
	} else {
		fmr.settings.Logger.Info("Body: No Body Supplied\n")
	}
	bodystring := prettyJSON.String()
	var result map[string]any
	json.Unmarshal([]byte(bodystring), &result)
	results := result["birds"].(map[string]any)

	for key, value := range results {
		// Each value is an `any` type, that is type asserted as a string
		fmt.Println(key, value.(string))
	}

	return bodystring

}
