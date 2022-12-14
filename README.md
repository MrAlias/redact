# Redact

[![Go Reference](https://pkg.go.dev/badge/github.com/MrAlias/redact.svg)](https://pkg.go.dev/github.com/MrAlias/redact)

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

Drop spans whose name is `"really-annoying-span"` or any that match `"health?check*"` (e.g. `"health-check"`, `"healthcheck-my-service"`).

```go
tracerProvider := trace.NewTracerProvider(
	redact.Span("really-annoying-span", "health?check*"),
	/* ... */
)
```

[OpenTelemetry]: https://opentelemetry.io/
[TracerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProvider
