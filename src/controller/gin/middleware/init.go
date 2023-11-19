package middleware

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/pecolynx/golang-webapi-boilerplate/src/controller/middleware")
