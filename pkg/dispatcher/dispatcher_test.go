/*
Copyright 2020 The Knative Authors

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

package dispatcher

import (
	"testing"
)

const (
	rabbitURL            = "amqp://localhost:5672/%2f"
	queueName            = "queue"
	exchangeName         = "default/knative-testbroker"
	eventData            = `{"testdata":"testdata"}`
	eventData2           = `{"testdata":"testdata2"}`
	responseData         = `{"testresponse":"testresponsedata"}`
	expectedData         = `"{\"testdata\":\"testdata\"}"`
	expectedData2        = `"{\"testdata\":\"testdata2\"}"`
	expectedResponseData = `"{\"testresponse\":\"testresponsedata\"}"`
)

func TestReadSpan(t *testing.T) {

}
