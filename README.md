# Redact

An [OpenTelemetry] [SpanProcessor] to redact tracing data.

## Getting Started

Wrap your existing [SpanProcessor] to redact data from tracing spans.

### Redact Attributes

Remove attributes from spans that have keys that match `"password"`, `"user"`, and `"secret"`.

```go
tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(
	redact.Attributes(YourSpanProcessor, "password", "user", "secret"),
))
```

### Redact Spans based on name

Drop spans whose name is `"health-check"`.

```go
tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(
	redact.Span(YourSpanProcessor, "health-check"),
))
```

### Redact Spans from an instrumentation scope

Drop spans from the `"noisy"` instrumentation library.

```go
tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(
	redact.Scope(YourSpanProcessor, instrumentation.Scope{Name: "noisy"}),
))
```

[OpenTelemetry]: https://opentelemetry.io/
[SpanProcessor]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#SpanProcessor
