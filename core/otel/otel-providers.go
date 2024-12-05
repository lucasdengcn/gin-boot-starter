package otel

import (
	"context"
	"errors"
	"gin-boot-starter/core/config"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var shutdownFuncs []func(context.Context) error

// Initialize a gRPC connection to be used by both the tracer and meter
// providers.
func initConn() *grpc.ClientConn {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	//
	otelCfg := config.GetConfig().OTEL
	//
	// secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	// if otelCfg.Insecure == "true" {
	// 	secureOption = otlptracegrpc.WithInsecure()
	// }
	//

	// Note the use of insecure transport here. TLS is recommended in production.
	credentials := insecure.NewCredentials()
	option := grpc.WithTransportCredentials(credentials)
	conn, err := grpc.NewClient(otelCfg.OTLPEndpoint, option)
	if err != nil {
		log.Fatalf("failed to create gRPC connection to collector: %v", err)
	}
	return conn
}

func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) func(context.Context) error {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown
}

// Initializes an OTLP exporter, and configures the corresponding meter provider.
func initMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) func(context.Context) error {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create metrics exporter: %v", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown
}

// Initializes an OTLP exporter, and configures the corresponding logging provider.
func initLogProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) func(context.Context) error {
	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create log exporter: %v", err)
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	global.SetLoggerProvider(loggerProvider)

	return loggerProvider.Shutdown
}

func InitProviders(ctx context.Context) {
	//
	otelCfg := config.GetConfig().OTEL
	if otelCfg.OTLPEndpoint == "" {
		return
	}
	//
	conn := initConn()

	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			attribute.String("service.name", otelCfg.ServiceName),
			attribute.String("service.version", otelCfg.ServiceName),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	if otelCfg.Tracer {
		shutdownTracerProvider := initTracerProvider(ctx, res, conn)
		shutdownFuncs = append(shutdownFuncs, shutdownTracerProvider)
	}
	//
	if otelCfg.Metric {
		shutdownMeterProvider := initMeterProvider(ctx, res, conn)
		shutdownFuncs = append(shutdownFuncs, shutdownMeterProvider)
	}
	//
	if otelCfg.Logging {
		shutdownLogProvider := initLogProvider(ctx, res, conn)
		shutdownFuncs = append(shutdownFuncs, shutdownLogProvider)
	}
}

func Shutdown(ctx context.Context) error {
	if len(shutdownFuncs) == 0 {
		return nil
	}
	var err error
	for _, fn := range shutdownFuncs {
		err = errors.Join(err, fn(ctx))
	}
	shutdownFuncs = nil
	if err != nil {
		log.Printf("OTEL Shutdown. Error: %v", err)
	} else {
		log.Printf("OTEL Shutdown.")
	}
	return err
}
