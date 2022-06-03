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

package utils

import (
	"context"

	"knative.dev/eventing-rabbitmq/vendor/knative.dev/pkg/logging"
	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
)

func SetBackoffPolicy(ctx context.Context, backoffPolicy string) eventingduckv1.BackoffPolicyType {
	if backoffPolicy == "" || backoffPolicy == "exponential" {
		return eventingduckv1.BackoffPolicyExponential
	} else if backoffPolicy == "linear" {
		return eventingduckv1.BackoffPolicyLinear
	}
	logging.FromContext(ctx).Fatalf("Invalid BACKOFF_POLICY specified: must be %q or %q", eventingduckv1.BackoffPolicyExponential, eventingduckv1.BackoffPolicyLinear)
	return ""
}
