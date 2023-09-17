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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	api "go.opentelemetry.io/otel/trace"
)

type attrRecorder struct {
	attrs []attribute.KeyValue
}

func (r *attrRecorder) OnEnd(s trace.ReadOnlySpan) {
	r.attrs = s.Attributes()
}
func (*attrRecorder) Shutdown(context.Context) error                   { return nil }
func (*attrRecorder) ForceFlush(context.Context) error                 { return nil }
func (*attrRecorder) OnStart(_ context.Context, _ trace.ReadWriteSpan) {}

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
	t.Run("EmptyAfterCreation", func(t *testing.T) {
		got := testAttributesAfterCreation(Attributes(), name, passStr, eID)
		contains(t, got, name, eID, passStr)
	})

	t.Run("SingleStringAttribute", func(t *testing.T) {
		got := testAttributes(Attributes(key), name, passStr, eID)
		contains(t, got, name, eID, replaced)
	})
	t.Run("SingleStringAttributeAfterCreation", func(t *testing.T) {
		got := testAttributesAfterCreation(Attributes(key), name, passStr, eID)
		contains(t, got, name, eID, replaced)
	})

	t.Run("NoMatchingKey", func(t *testing.T) {
		got := testAttributes(Attributes("secret"), name, passStr, eID)
		contains(t, got, name, eID, passStr)
	})
	t.Run("NoMatchingKeyAfterCreation", func(t *testing.T) {
		got := testAttributesAfterCreation(Attributes("secret"), name, passStr, eID)
		contains(t, got, name, eID, passStr)
	})

	t.Run("DifferentValueTypes", func(t *testing.T) {
		got := testAttributes(Attributes(key), name, passBool, eID)
		contains(t, got, name, eID, replaced)
	})
	t.Run("DifferentValueTypesAfterCreation", func(t *testing.T) {
		got := testAttributesAfterCreation(Attributes(key), name, passBool, eID)
		contains(t, got, name, eID, replaced)
	})
}

func testAttributes(opt trace.TracerProviderOption, attrs ...attribute.KeyValue) []attribute.KeyValue {
	r := &attrRecorder{}
	tp := trace.NewTracerProvider(opt, trace.WithSpanProcessor(r))
	defer func() { _ = tp.Shutdown(context.Background()) }()

	ctx := context.Background()
	tracer := tp.Tracer("testAttributes")
	_, s := tracer.Start(ctx, "span name", api.WithAttributes(attrs...))
	s.End()
	return r.attrs
}

func testAttributesAfterCreation(opt trace.TracerProviderOption, attrs ...attribute.KeyValue) []attribute.KeyValue {
	r := &attrRecorder{}
	tp := trace.NewTracerProvider(opt, trace.WithSpanProcessor(r))
	defer func() { _ = tp.Shutdown(context.Background()) }()

	ctx := context.Background()
	tracer := tp.Tracer("testAttributes")
	_, s := tracer.Start(ctx, "span name")
	s.SetAttributes(attrs...)
	s.End()
	return r.attrs
}

func BenchmarkAttributeCensorOnStart(b *testing.B) {
	b.Run("0/16", benchAttributeCensorOnStart(0, 16))
	b.Run("1/16", benchAttributeCensorOnStart(1, 16))
	b.Run("2/16", benchAttributeCensorOnStart(2, 16))
	b.Run("4/16", benchAttributeCensorOnStart(4, 16))
	b.Run("8/16", benchAttributeCensorOnStart(8, 16))
	b.Run("16/16", benchAttributeCensorOnStart(16, 16))
}

type rwSpan struct {
	trace.ReadWriteSpan

	attrs []attribute.KeyValue
}

func (rwSpan) SetAttributes(...attribute.KeyValue) {}
func (s rwSpan) Attributes() []attribute.KeyValue {
	return s.attrs
}

func benchAttributeCensorOnStart(redacted, total int) func(*testing.B) {
	if redacted > total {
		panic("redacted needs to be less than or equal to total")
	}

	replacements := make(map[attribute.Key]attribute.Value)
	attrs := make([]attribute.KeyValue, total)
	for i := range attrs {
		key := attribute.Key(strconv.Itoa(i))
		if i < redacted {
			replacements[key] = attribute.StringValue(defaultReplace)
		}
		attrs[i] = attribute.KeyValue{
			Key:   key,
			Value: attribute.IntValue(i),
		}
	}

	s := rwSpan{attrs: attrs}
	ac := NewAttributeCensor(replacements)
	ctx := context.Background()

	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ac.OnStart(ctx, s)
			ac.OnEnd(s)
		}
	}
}
