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
		passStr  = attribute.String(key, "super-secret-pswd")
		passBool = attribute.Bool(key, true)
		replaced = attribute.String(key, defaultReplace)
	)

	r := &attrRecorder{}
	tp := trace.NewTracerProvider(Attributes(key), trace.WithSpanProcessor(r))
	t.Cleanup(func() { tp.Shutdown(context.Background()) })

	tracer := tp.Tracer("TestAttributes")
	tracer.Start(context.Background(), "span name", api.WithAttributes(passStr, name))
	assert.Contains(t, r.attrs, name)
	assert.Contains(t, r.attrs, replaced)

	tracer.Start(context.Background(), "span name", api.WithAttributes(name, passBool))
	assert.Contains(t, r.attrs, name)
	assert.Contains(t, r.attrs, replaced)
}
