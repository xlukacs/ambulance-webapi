package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/technologize/otel-go-contrib/otelginmetrics"
	"github.com/xlukacs/ambulance-webapi/api"
	"github.com/xlukacs/ambulance-webapi/internal/ambulance_wl"
	"github.com/xlukacs/ambulance-webapi/internal/db_service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup telemetry
	shutdown, err := initTelemetry()
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()
	engine.Use(
		otelginmetrics.Middleware(
			"Ambulance WebAPI Service",
			// Custom attributes
			otelginmetrics.WithAttributes(func(serverName, route string, request *http.Request) []attribute.KeyValue {
				return append(otelginmetrics.DefaultAttributes(serverName, route, request))
			}),
		),
		// otelgin.Middleware(serverName), TODO this needs to be here, but where is the serverName coming from...???
	)

	// setup context update  middleware
	dbService := db_service.NewMongoService[ambulance_wl.Ambulance](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	// request routings
	ambulance_wl.AddRoutes(engine)
	engine.GET("/openapi", api.HandleOpenApi)

	// metrics endpoint
	promhandler := promhttp.Handler()
	engine.Any("/metrics", func(ctx *gin.Context) {
		promhandler.ServeHTTP(ctx.Writer, ctx.Request)
	})
	engine.Run(":" + port)
}

// initialize OpenTelemetry instrumentations
func initTelemetry() (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String("Ambulance WebAPI Service")),
		resource.WithAttributes(semconv.ServiceNamespaceKey.String("WAC Hospital")),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithContainer(),
	)

	if err != nil {
		return nil, err
	}

	metricExporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	metricProvider := metric.NewMeterProvider(metric.WithReader(metricExporter), metric.WithResource(res))
	otel.SetMeterProvider(metricProvider)
	// setup trace exporter, only otlp supported
	// see also https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/exporters/autoexport
	traceExportType := os.Getenv("OTEL_TRACES_EXPORTER")
	if traceExportType == "otlp" {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		// we will configure exporter by using env variables defined
		// at https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/
		traceExporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			return nil, err
		}

		traceProvider := trace.NewTracerProvider(
			trace.WithResource(res),
			trace.WithSyncer(traceExporter))

		otel.SetTracerProvider(traceProvider)
		otel.SetTextMapPropagator(propagation.TraceContext{})
		// Shutdown function will flush any remaining spans
		return traceProvider.Shutdown, nil
	} else {
		// no otlp trace exporter configured
		noopShutdown := func(context.Context) error { return nil }
		return noopShutdown, nil
	}
}
