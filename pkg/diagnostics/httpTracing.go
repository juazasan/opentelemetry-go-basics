package diagnostics

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func HTTPTraceMiddleware(next http.Handler) http.Handler {
	tracer = otel.Tracer("httpTracing")
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		path := string(r.URL.Path)
		fmt.Println(fmt.Sprintf("%s %s", r.Method, path))
		_, span := startTracingClientSpanFromHTTPContext(r, path)
		next.ServeHTTP(rw, r)
		(*span).End()
	})
}

func startTracingClientSpanFromHTTPContext(request *http.Request, spanName string) (*context.Context, *trace.Span) {
	cxt := request.Context()
	//sc := SpanContextFromRequest(request)
	kindOption := trace.WithSpanKind(trace.SpanKindServer)
	ctx, span := tracer.Start(cxt, spanName, kindOption)
	return &ctx, &span
}

func SpanFromRequest(request *http.Request) trace.Span {
	sc := trace.SpanFromContext(request.Context())
	return sc
}
