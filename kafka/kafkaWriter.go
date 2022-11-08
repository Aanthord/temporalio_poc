package kafkaWriter

import (
	"context"
	"os"

	kafka "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	service     = "temporalio-KafkaWriter"
	environment = "test"
	id          = 1
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func main(topic string, key1 string, val1 string) {
	// get kafka writer using environment variables.
	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("KafkaURL")),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    os.Getenv("BatchSize"),
		BatchTimeout: os.Getenv("BatchTimeoutMS"),
	}

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key1),
			Value: []byte(val1),
		},
	)

	w.Close()
}
