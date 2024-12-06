package otel

import (
	"context"
	"errors"
	"gin-boot-starter/core/config"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	stdoutlog "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	stdoutmetric "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	stdouttrace "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

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
func initConn() (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	//
	otelCfg := config.GetConfig().OTEL
	if otelCfg.Stdout {
		return nil, nil
	}
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
		log.Error().Err(err).Msgf("failed to create gRPC connection to collector")
		return nil, err
	}
	return conn, nil
}

func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (*sdktrace.TracerProvider, error) {
	// Set up a trace exporter
	stdoutEnabled := config.GetConfig().OTEL.Stdout
	var exporterOption sdktrace.TracerProviderOption
	if stdoutEnabled {
		exporter, err := stdouttrace.New()
		if err != nil {
			return nil, err
		}
		exporterOption = sdktrace.WithBatcher(exporter)
	} else {
		exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			log.Error().Err(err).Msgf("failed to create trace exporter")
			return nil, err
		}
		exporterOption = sdktrace.WithBatcher(exporter)
	}
	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		exporterOption,
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider, nil
}

// Initializes an OTLP exporter, and configures the corresponding meter provider.
func initMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (*sdkmetric.MeterProvider, error) {
	stdoutEnabled := config.GetConfig().OTEL.Stdout
	var reader *sdkmetric.PeriodicReader
	if stdoutEnabled {
		// different type of exporter
		exporter, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		reader = sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second))
	} else {
		exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
		if err != nil {
			log.Error().Err(err).Msgf("failed to create metrics exporter")
			return nil, err
		}
		reader = sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second))
	}
	//
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}

// Initializes an OTLP exporter, and configures the corresponding logging provider.
func initLogProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (*sdklog.LoggerProvider, error) {
	stdoutEnabled := config.GetConfig().OTEL.Stdout
	var logProcessor *sdklog.BatchProcessor
	if stdoutEnabled {
		logExporter, err := stdoutlog.New()
		if err != nil {
			return nil, err
		}
		logProcessor = sdklog.NewBatchProcessor(logExporter, sdklog.WithExportInterval(10*time.Second))
	} else {
		logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
		if err != nil {
			log.Error().Err(err).Msgf("failed to create log exporter")
			return nil, err
		}
		logProcessor = sdklog.NewBatchProcessor(logExporter, sdklog.WithExportInterval(10*time.Second))
	}
	//
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(logProcessor),
		sdklog.WithResource(res),
	)
	//
	global.SetLoggerProvider(loggerProvider)
	//
	return loggerProvider, nil
}

func InitProviders(ctx context.Context) {
	//
	otelCfg := config.GetConfig().OTEL
	//
	conn, err := initConn()
	if err != nil {
		return
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			attribute.String("service.name", otelCfg.ServiceName),
			attribute.String("service.version", otelCfg.ServiceName),
		),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Resource create error")
	}

	if otelCfg.Tracer {
		tracerProvider, err := initTracerProvider(ctx, res, conn)
		if err != nil {
			return
		}
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	}
	//
	if otelCfg.Metric {
		meterProvider, err := initMeterProvider(ctx, res, conn)
		if err != nil {
			return
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	}
	//
	if otelCfg.Logging {
		logProvider, err := initLogProvider(ctx, res, conn)
		if err != nil {
			return
		}
		shutdownFuncs = append(shutdownFuncs, logProvider.Shutdown)
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
