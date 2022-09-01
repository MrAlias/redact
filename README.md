# Redact

An set of [OpenTelemetry] [TracerProviderOption]s to redact tracing data.

## Getting Started

Pass your needed redact option to a new [OpenTelemetry] [TracerProvider].

### Redact Attributes

Remove attributes from spans that have keys that match `"password"`, `"user"`, and `"secret"`.

```go
tracerProvider := trace.NewTracerProvider(
	redact.Attributes("password", "user", "secret"),
	/* ... */
)
```

### Redact Spans based on name

Drop spans whose name is `"health-check"`.

```go
tracerProvider := trace.NewTracerProvider(
	redact.Span("health-check"),
	/* ... */
)
```

### Redact Spans from an instrumentation scope

Drop spans from the `"noisy"` instrumentation library.

```go
tracerProvider := trace.NewTracerProvider(
	redact.Scope(instrumentation.Scope{Name: "noisy"}),
	/* ... */
)
```

[OpenTelemetry]: https://opentelemetry.io/
[TracerProviderOption]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProviderOption
[TracerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProvider
