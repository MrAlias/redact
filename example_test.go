// Copyright 2022 Tyler Yahn (MrAlias)
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

package redact_test

import (
	"github.com/MrAlias/redact"
	"go.opentelemetry.io/otel/sdk/trace"
)

func ExampleAttributes() {
	_ = trace.NewTracerProvider(
		// Replace attributes from new spans that have keys matching
		// "password", "user", and "secret" with the redacted value
		// "****REDACTED****".
		redact.Attributes("password", "user", "secret"),
		/* From here, configure your trace pipeline as normal ... */
	)
}

func ExampleSpan() {
	_ = trace.NewTracerProvider(
		// Drop spans whose name is `"health-check"`.
		redact.Span("health-check"),
		/* From here, configure your trace pipeline as normal ... */
	)
}
