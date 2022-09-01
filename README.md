# Redact

An set of [OpenTelemetry] [TracerProviderOption]s to redact tracing data.

## Getting Started

Pass your needed redact option to a new [OpenTelemetry] [TracerProvider].

### Redact Attributes

Replace attributes from new spans that have keys matching `"password"`, `"user"`, and `"secret"` with a redacted value.

```go
tracerProvider := trace.NewTracerProvider(
	redact.Attributes("password", "user", "secret"),
	/* ... */
)
```

### (TODO) Redact Spans based on name

Drop spans whose name is `"health-check"`.

```go
tracerProvider := trace.NewTracerProvider(
	redact.Span("health-check"),
	/* ... */
)
```

### (TODO) Redact Spans from an instrumentation scope

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
