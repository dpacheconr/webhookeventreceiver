// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhookeventreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/webhookeventreceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type webhookeventreceiverReceiver struct {
	host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Logs
	config       *Config
}

func (webhookeventreceiverRcvr *webhookeventreceiverReceiver) Start(ctx context.Context, host component.Host) error {
	webhookeventreceiverRcvr.host = host
	// ctx = context.Background()
	// ctx, webhookeventreceiverRcvr.cancel = context.WithCancel(ctx)
	webhookeventreceiverRcvr.logger.Info("I should start processing now!")
	return nil
}

func (webhookeventreceiverRcvr *webhookeventreceiverReceiver) Shutdown(ctx context.Context) error {
	webhookeventreceiverRcvr.cancel()
	return nil
}
