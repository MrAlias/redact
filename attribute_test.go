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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	api "go.opentelemetry.io/otel/trace"
)

type attrRecorder struct {
	attrs []attribute.KeyValue
}

func (*attrRecorder) OnEnd(trace.ReadOnlySpan)         {}
func (*attrRecorder) Shutdown(context.Context) error   { return nil }
func (*attrRecorder) ForceFlush(context.Context) error { return nil }
func (r *attrRecorder) OnStart(_ context.Context, s trace.ReadWriteSpan) {
	r.attrs = s.Attributes()
}

func TestAttributes(t *testing.T) {
	const key = "password"
	var (
		name     = attribute.String("name", "bob")
		eID      = attribute.Int("employee-id", 9287)
		passStr  = attribute.String(key, "super-secret-pswd")
		passBool = attribute.Bool(key, true)
		replaced = attribute.String(key, defaultReplace)
	)

	contains := func(t *testing.T, got []attribute.KeyValue, want ...attribute.KeyValue) {
		t.Helper()
		for _, w := range want {
			assert.Contains(t, got, w)
		}
	}

	t.Run("Empty", func(t *testing.T) {
		got := testAttributes(Attributes(), name, passStr, eID)
		contains(t, got, name, eID, passStr)
	})

	t.Run("SingleStringAttribute", func(t *testing.T) {
		got := testAttributes(Attributes(key), name, passStr, eID)
		contains(t, got, name, eID, replaced)
	})

	t.Run("NoMatchingKey", func(t *testing.T) {
		got := testAttributes(Attributes("secret"), name, passStr, eID)
		contains(t, got, name, eID, passStr)
	})

	t.Run("DifferentValueTypes", func(t *testing.T) {
		got := testAttributes(Attributes(key), name, passBool, eID)
		contains(t, got, name, eID, replaced)
	})
}

func testAttributes(opt trace.TracerProviderOption, attrs ...attribute.KeyValue) []attribute.KeyValue {
	r := &attrRecorder{}
	tp := trace.NewTracerProvider(opt, trace.WithSpanProcessor(r))
	defer func() { tp.Shutdown(context.Background()) }()

	ctx := context.Background()
	tracer := tp.Tracer("testAttributes")
	_, s := tracer.Start(ctx, "span name", api.WithAttributes(attrs...))
	s.End()
	return r.attrs
}

func BenchmarkAttribute(b *testing.B) {
	b.Run("Keys/0", benchmarkAttribute())
	b.Run("Keys/1", benchmarkAttribute("secret"))
	b.Run("Keys/2", benchmarkAttribute("secret", "password"))
}

func benchmarkAttribute(keys ...string) func(*testing.B) {
	ctx := context.Background()
	tp := trace.NewTracerProvider(Attributes(keys...))

	tracer := tp.Tracer("BenchmarkAttribute")
	attrs := []attribute.KeyValue{
		attribute.Bool("bool", true),
		attribute.Int("secret", 42),
		attribute.String("password", "secret password"),
	}

	return func(b *testing.B) {
		b.Cleanup(func() { tp.Shutdown(ctx) })
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, s := tracer.Start(ctx, "span name", api.WithAttributes(attrs...))
			s.End()
		}
	}
}
