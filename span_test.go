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
)

func TestSpan(t *testing.T) {
	names := []string{"health-check", "admin"}
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(
		Span(names...),
		trace.WithSpanProcessor(sr),
	)

	ctx := context.Background()
	tracer := tp.Tracer("TestSpan")
	for _, n := range names {
		_, s := tracer.Start(ctx, n)
		assert.False(t, s.IsRecording(), "span should not be recorded")
		assert.False(t, s.SpanContext().IsSampled(), "span should not be sampled")
		s.End()
	}

	const valid = "HTTP GET"
	_, s := tracer.Start(ctx, valid)
	assert.True(t, s.IsRecording())
	assert.True(t, s.SpanContext().IsSampled())
	s.End()

	require.NoError(t, tp.Shutdown(ctx))

	got := sr.Ended()
	require.Len(t, got, 1, "only one span should be exported")
	assert.Equal(t, valid, got[0].Name())
}

func TestSpanCensorDescription(t *testing.T) {
	sc := NewSpanCensor(trace.AlwaysSample())
	expect := fmt.Sprintf("SpanCensor(%s)", trace.AlwaysSample().Description())
	assert.Equal(t, expect, sc.Description())
}
