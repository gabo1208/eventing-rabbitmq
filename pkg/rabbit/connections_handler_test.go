/*
Copyright 2022 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rabbit

import (
	"errors"
	"testing"

	"go.uber.org/zap"
)

func Test_ValidSetupRabbitMQ(t *testing.T) {
	logger := zap.NewNop().Sugar()
	rabbitMQHelper := NewRabbitMQHelper(1, logger, ValidDial).(*RabbitMQHelper)
	rabbitMQHelper.SetupRabbitMQConnectionAndChannel("amqp://localhost:5672/%2f", ConfigTest)
	if rabbitMQHelper.Connection == nil || rabbitMQHelper.Channel == nil {
		t.Errorf("rabbitMQHelper connection and channel should be set %s %s", rabbitMQHelper.Connection, rabbitMQHelper.Channel)
	}
	rabbitMQHelper.CloseRabbitMQConnections()
}

func Test_InvalidSetupRabbitMQ(t *testing.T) {
	var err error
	retryConnection, retryChannel := true, true
	logger := zap.NewNop().Sugar()
	// test invalid connection setup
	rabbitMQHelper := NewRabbitMQHelper(1, logger, BadConnectionDial).(*RabbitMQHelper)
	rabbitMQHelper.Connection, rabbitMQHelper.Channel, err = rabbitMQHelper.CreateAndConfigConnectionsAndChannel(&retryConnection, &retryChannel, "amqp://localhost:5672/%2f", nil)
	if err == nil || !retryConnection || rabbitMQHelper.GetConnection() != nil {
		t.Errorf("unexpected error == nil when setting up invalid connection %s %s %s", rabbitMQHelper.GetConnection(), rabbitMQHelper.GetChannel(), err)
	}
	rabbitMQHelper.Connection, rabbitMQHelper.Channel, err = nil, nil, nil
	rabbitMQHelper.CloseRabbitMQConnections()
	// test invalid channel setup
	rabbitMQHelper = NewRabbitMQHelper(1, logger, BadChannelDial).(*RabbitMQHelper)
	rabbitMQHelper.Connection, rabbitMQHelper.Channel, err = rabbitMQHelper.CreateAndConfigConnectionsAndChannel(&retryConnection, &retryChannel, "amqp://localhost:5672/%2f", ConfigTest)
	if err == nil || !retryChannel || rabbitMQHelper.GetChannel() != nil {
		t.Errorf("unexpected error == nil when setting up invalid channel %s %s %s", rabbitMQHelper.GetConnection(), rabbitMQHelper.GetChannel(), err)
	}
	rabbitMQHelper.CloseRabbitMQConnections()
	rabbitMQHelper.Connection, rabbitMQHelper.Channel, err = nil, nil, nil
	retryConnection, retryChannel = true, true
	// test invalid config setup
	rabbitMQHelper = NewRabbitMQHelper(1, logger, ValidDial).(*RabbitMQHelper)
	rabbitMQHelper.Connection, rabbitMQHelper.Channel, err = rabbitMQHelper.CreateAndConfigConnectionsAndChannel(&retryConnection, &retryChannel, "amqp://localhost:5672/%2f", InvalidConfigTest)
	if err == nil || !retryConnection || !retryChannel || rabbitMQHelper.GetConnection() == nil || rabbitMQHelper.GetChannel() == nil {
		t.Errorf("unexpected error == nil when setting up invalid config %s %s %s", rabbitMQHelper.GetConnection(), rabbitMQHelper.GetChannel(), err)
	}
	rabbitMQHelper.CloseRabbitMQConnections()
}

func ConfigTest(conn RabbitMQConnectionInterface, channel RabbitMQChannelInterface) error {
	ChannelConfirm(conn, channel)
	ChannelQoS(conn, channel)
	return nil
}

func InvalidConfigTest(conn RabbitMQConnectionInterface, channel RabbitMQChannelInterface) error {
	return errors.New("invalid config test")
}

func TestAdapter_VhostHandler(t *testing.T) {
	for _, tt := range []struct {
		name   string
		broker string
		vhost  string
		want   string
	}{{
		name:   "no broker nor vhost",
		broker: "",
		vhost:  "",
		want:   "",
	}, {
		name:   "no vhost",
		broker: "amqp://localhost:5672",
		vhost:  "",
		want:   "amqp://localhost:5672",
	}, {
		name:   "no broker",
		broker: "",
		vhost:  "test-vhost",
		want:   "test-vhost",
	}, {
		name:   "broker and vhost without separating slash",
		broker: "amqp://localhost:5672",
		vhost:  "test-vhost",
		want:   "amqp://localhost:5672/test-vhost",
	}, {
		name:   "broker and vhost without separating slash but vhost with ending slash",
		broker: "amqp://localhost:5672",
		vhost:  "test-vhost/",
		want:   "amqp://localhost:5672/test-vhost/",
	}, {
		name:   "broker with trailing slash and vhost without the slash",
		broker: "amqp://localhost:5672/",
		vhost:  "test-vhost",
		want:   "amqp://localhost:5672/test-vhost",
	}, {
		name:   "vhost starting with slash and broker without the slash",
		broker: "amqp://localhost:5672",
		vhost:  "/test-vhost",
		want:   "amqp://localhost:5672/test-vhost",
	}, {
		name:   "broker and vhost with slash",
		broker: "amqp://localhost:5672/",
		vhost:  "/test-vhost",
		want:   "amqp://localhost:5672//test-vhost",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			got := VHostHandler(tt.broker, tt.vhost)
			if got != tt.want {
				t.Errorf("Unexpected URI for %s/%s want:\n%+s\ngot:\n%+s", tt.broker, tt.vhost, tt.want, got)
			}
		})
	}
}
