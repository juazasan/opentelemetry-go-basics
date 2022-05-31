package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

type Client struct {
	serverURI string
	tp        *sdktrace.TracerProvider
}

func NewClient(serverURI string) Client {
	tracer = otel.Tracer("client")
	exp, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		panic(err)
	}
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("client"),
		semconv.ServiceVersionKey.String("1.0.0"),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return Client{serverURI, tp}
}

func (c *Client) Run(ctx context.Context) {

	for i := 0; i < 5; i++ {
		client := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
		client.Timeout = 1 * time.Second
		ctx, span := tracer.Start(ctx, "get:Hello")
		req, err := http.NewRequestWithContext(ctx, "GET", c.serverURI, nil)
		if err != nil {
			panic(err)
		}

		defer span.End()
		res, err := client.Do(req)
		if err != nil {
			panic(nil)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(nil)
		}
		fmt.Println(string(body))
	}
}

func (c Client) Shutdown() {
	ctx := context.Background()
	_ = c.tp.Shutdown(ctx)
}
