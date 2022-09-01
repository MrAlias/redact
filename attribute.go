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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

const defaultReplace = "****REDACTED****"

// Attributes returns an OpenTelemetry SDK TracerProviderOption. It registers
// an OpenTelemetry SpanProcessor that redacts attributes of new spans matching
// the passed keys.
func Attributes(keys ...string) trace.TracerProviderOption {
	r := make(map[attribute.Key]attribute.Value, len(keys))
	for _, k := range keys {
		r[attribute.Key(k)] = attribute.StringValue(defaultReplace)
	}
	censor := NewAttributeCensor(r)
	return trace.WithSpanProcessor(censor)
}

// AttributeCensor is an OpenTelemetry SpanProcessor that censors attributes of
// new spans.
type AttributeCensor struct {
	// args is a slice allocated on creation that is reused when calling
	// SetAttributes in OnStart.
	args         []attribute.KeyValue
	replacements map[attribute.Key]attribute.Value
}

// NewAttributeCensor returns an AttributeCensor that uses the provided mapping
// of replacement values for a set of keys to redact matching attributes.
// Attributes are matched based on the equality of keys.
func NewAttributeCensor(replacements map[attribute.Key]attribute.Value) AttributeCensor {
	a := AttributeCensor{
		// Allocate a reusable slice to pass to SetAttributes.
		args:         make([]attribute.KeyValue, 0, len(replacements)),
		replacements: replacements,
	}
	return a
}

// OnStart censors the attributes of s matching the Replacements keys of c.
func (c AttributeCensor) OnStart(_ context.Context, s trace.ReadWriteSpan) {
	// The SetAttributes method only overwrites attributes that already exists,
	// it does not set the attributes to only the values passed. Therefore,
	// determine if there are any attributes that need to be redacted and set
	// overrides for only those values.
	redacted := c.args[:0]
	for _, a := range s.Attributes() {
		if v, ok := c.replacements[a.Key]; ok {
			redacted = append(redacted, attribute.KeyValue{
				Key:   a.Key,
				Value: v,
			})
		}
	}
	if len(redacted) > 0 {
		s.SetAttributes(redacted...)
	}
}

// OnEnd does nothing.
func (AttributeCensor) OnEnd(trace.ReadOnlySpan) {}

// Shutdown does nothing.
func (AttributeCensor) Shutdown(context.Context) error { return nil }

// ForceFlush does nothing.
func (AttributeCensor) ForceFlush(context.Context) error { return nil }
