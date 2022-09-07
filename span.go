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
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/sdk/trace"
	api "go.opentelemetry.io/otel/trace"
)

// Span returns an OpenTelemetry SDK TracerProviderOption. It registers an
// OpenTelemetry Sampler that drops all spans with matching names. A name can
// be the exact span name, or contain wildcards. See NewSpanCensor for
// information about matching.
func Span(names ...string) trace.TracerProviderOption {
	parent := trace.ParentBased(trace.AlwaysSample())
	return trace.WithSampler(NewSpanCensor(parent, names...))
}

// SpanCensor is an OpenTelemetry Sampler that drops spans with certain names.
// See NewSpanCensor for information about name matching.
type SpanCensor struct {
	wrapped trace.Sampler
	desc    string
	exact   map[string]struct{}
	wcRe    *regexp.Regexp
}

// NewSpanCensor returns a new configured SpanCensor. Any span with a name
// matching one of the passed names will be dropped. A name can be an exact
// string match for the span name, or contain wildcards. The * wildcard is
// expanded to match zero or more characters and the ? wildcard matches exactly
// one character. One or more wildcards can be used in combination to match a
// name.
func NewSpanCensor(parent trace.Sampler, names ...string) SpanCensor {
	var wc []string
	n := make(map[string]struct{}, len(names))
	for _, name := range names {
		if strings.ContainsAny(name, "*?") {
			name = regexp.QuoteMeta(name)
			name = strings.ReplaceAll(name, "\\?", ".")
			name = strings.ReplaceAll(name, "\\*", ".*")
			wc = append(wc, name)
		} else {
			n[name] = struct{}{}
		}
	}

	var wcRe *regexp.Regexp
	if len(wc) > 0 {
		wcRe = regexp.MustCompile("^(?:" + strings.Join(wc, "|") + ")$")
	}
	return SpanCensor{
		wrapped: parent,
		desc:    fmt.Sprintf("SpanCensor(%s)", parent.Description()),
		exact:   n,
		wcRe:    wcRe,
	}
}

func (s SpanCensor) match(name string) bool {
	_, match := s.exact[name]
	return match || (s.wcRe != nil && s.wcRe.MatchString(name))
}

// ShouldSample returns a sampling decision to drop when p contains a name
// matching the SpanCensor criteria. Otherwise, the sampling decision is
// delegated to the Sampler SpanCensor wraps.
func (s SpanCensor) ShouldSample(p trace.SamplingParameters) trace.SamplingResult {
	if s.match(p.Name) {
		return trace.SamplingResult{
			Decision:   trace.Drop,
			Tracestate: api.SpanContextFromContext(p.ParentContext).TraceState(),
		}
	}
	return s.wrapped.ShouldSample(p)
}

// Description returns an identification description of the SpanCensor.
func (s SpanCensor) Description() string { return s.desc }
