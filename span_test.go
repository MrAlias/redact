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
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	api "go.opentelemetry.io/otel/trace"
)

func TestSpan(t *testing.T) {
	tester := NewSpanCensorTester(
		"admin",
		"HTTP *",
		"health?check*",
		"client-??-op",
	)

	tests := []struct {
		name     string
		redacted bool
	}{
		{name: "admin", redacted: true},
		{name: "HTTP GET", redacted: true},
		{name: "HTTP POST", redacted: true},
		{name: "health-check", redacted: true},
		{name: "health-check-srv", redacted: true},
		{name: "health_check my-service", redacted: true},
		{name: "client-00-op", redacted: true},
		{name: "client-ab-op", redacted: true},
		{name: "client-1-op", redacted: false},
		{name: "DB GET", redacted: false},
		{name: "RPC serveData", redacted: false},
	}
	for _, test := range tests {
		if test.redacted {
			t.Run("Redacted/"+test.name, tester.RunRedacted(test.name))
		} else {
			t.Run("Valid/"+test.name, tester.RunValid(test.name))
		}
	}
}

type SpanCensorTester struct {
	names []string
}

func NewSpanCensorTester(names ...string) *SpanCensorTester {
	return &SpanCensorTester{names: names}
}

func (sct *SpanCensorTester) run(name string) (*tracetest.SpanRecorder, *trace.TracerProvider, api.Span) {
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(
		Span(sct.names...),
		trace.WithSpanProcessor(sr),
	)
	_, s := tp.Tracer("SpanCensorTest").Start(context.Background(), name)
	return sr, tp, s
}

func (sct *SpanCensorTester) RunRedacted(name string) func(*testing.T) {
	sr, tp, s := sct.run(name)
	return func(t *testing.T) {
		assert.Falsef(t, s.IsRecording(), "%q recorded", name)
		assert.Falsef(t, s.SpanContext().IsSampled(), "%q sampled", name)
		s.End()

		require.NoError(t, tp.Shutdown(context.Background()))
		got := sr.Ended()
		assert.Lenf(t, got, 0, "%q exported", name)
	}
}

func (sct *SpanCensorTester) RunValid(name string) func(*testing.T) {
	sr, tp, s := sct.run(name)
	return func(t *testing.T) {
		assert.Truef(t, s.IsRecording(), "%q not recorded", name)
		assert.Truef(t, s.SpanContext().IsSampled(), "%q not sampled", name)
		s.End()

		require.NoError(t, tp.Shutdown(context.Background()))
		got := sr.Ended()
		require.Len(t, got, 1, "only one span should be exported")
		assert.Equal(t, name, got[0].Name(), "exported wrong span")
	}
}

func TestSpanCensorDescription(t *testing.T) {
	sc := NewSpanCensor(trace.AlwaysSample())
	expect := fmt.Sprintf("SpanCensor(%s)", trace.AlwaysSample().Description())
	assert.Equal(t, expect, sc.Description())
}
