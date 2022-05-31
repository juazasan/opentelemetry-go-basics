package runtime

import (
	"context"

	"github.com/juazasan/opentelemetry-go-basics/pkg/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type myRuntime struct {
	//ctx    context.Context
	server server.Server
	tp     *sdktrace.TracerProvider
	//runningServers []io.Closer
}

func NewRuntime() myRuntime {

	exp, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		panic(err)
	}
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("server"),
		semconv.ServiceVersionKey.String("1.0.0"),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	server := server.NewServer()
	return myRuntime{server, tp}
}

func (a *myRuntime) Run() {
	a.startHTTPServer()
}

func (a *myRuntime) startHTTPServer() {
	a.server.StartNotBlocking()
}

func (a *myRuntime) Shutdown() {
	ctx := context.Background()
	_ = a.tp.Shutdown(ctx)
}
