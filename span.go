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

package redact // import "github.com/MrAlias/redact"

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/trace"
	api "go.opentelemetry.io/otel/trace"
)

func Span(names ...string) trace.TracerProviderOption {
	parent := trace.ParentBased(trace.AlwaysSample())
	return trace.WithSampler(NewSpanCensor(parent, names...))
}

type SpanCensor struct {
	wrapped trace.Sampler
	desc    string
	names   map[string]struct{}
}

func NewSpanCensor(parent trace.Sampler, names ...string) SpanCensor {
	n := make(map[string]struct{}, len(names))
	for _, name := range names {
		n[name] = struct{}{}
	}
	return SpanCensor{
		wrapped: parent,
		desc:    fmt.Sprintf("SpanCensor(%s)", parent.Description()),
		names:   n,
	}
}

func (s SpanCensor) ShouldSample(p trace.SamplingParameters) trace.SamplingResult {
	if _, drop := s.names[p.Name]; drop {
		return trace.SamplingResult{
			Decision:   trace.Drop,
			Tracestate: api.SpanContextFromContext(p.ParentContext).TraceState(),
		}
	}
	return s.wrapped.ShouldSample(p)
}

func (s SpanCensor) Description() string { return s.desc }
