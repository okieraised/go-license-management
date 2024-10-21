package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"time"
)

var tp *sdktrace.TracerProvider
var tExp *otlptrace.Exporter

func GetInstance() *sdktrace.TracerProvider {
	return tp
}

func NewTracerProvider(grpcHost string, serviceName, namespace string) error {
	var err error
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientConn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(clientConn),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{PermitWithoutStream: true})),
	)

	tExp, err = otlptrace.New(timeoutCtx, traceClient)
	if err != nil {
		return err
	}

	res, err := resource.New(timeoutCtx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.K8SNamespaceName(namespace),
		),
	)
	if err != nil {
		return err
	}
	bsp := sdktrace.NewBatchSpanProcessor(tExp)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)
	return nil
}

func Shutdown() {
	if tExp == nil {
		return
	}
	ctx := context.Background()
	cCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	err := tExp.Shutdown(cCtx)
	if err != nil {
		otel.Handle(err)
	}
}
