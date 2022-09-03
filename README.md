# Redact

Unfortunately, you are here.
You have [OpenTelemetry] tracing data that shouldn't exist and you want it gone.
Ideally, you are able to stop the generation of this data.
But until that day arrives, `redact` can help!

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

### Redact Spans based on name

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
[TracerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProvider
