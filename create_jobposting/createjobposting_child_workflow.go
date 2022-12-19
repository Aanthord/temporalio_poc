package createjobposting_child_workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aanthord/temporalio_poc/legacy"
	"github.com/aanthord/temporalio_poc/watson"
	kafka "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	service     = "temporalio-createjobposting"
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

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func CreateJobPostingChildWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	// The client is a heavyweight object that should be created only once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "child-workflow", worker.Options{})

	w.RegisterWorkflow(child_workflow.CreateJobPostingParentWorkflow)
	w.RegisterWorkflow(child_workflow.CreateJobPostingChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)

		// get kafka reader using environment variables.
		kafkaURL := os.Getenv("kafkaURL")
		topic := os.Getenv("topic")
		groupID := os.Getenv("groupID")

		reader := getKafkaReader(kafkaURL, topic, groupID)

		defer reader.Close()

		fmt.Println("start consuming ... !!")
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			logger.Info("Consuming message")
			var payload interface{} // The interface where we will save the converted JSON data.

			json.Unmarshal(m, &payload)           // Convert JSON data into interface{} type
			m := payload.(map[string]interface{}) // To use the converted data we will need to convert it
			// into a map[string]interface{}
			logger.Info("Getting user_id")
			profile := legacy.LegacyGet(string(m.Userid))

			err := os.WriteFile("/tmp/"+string(m.Userid)+".json", profile, 0644)
			check(err)

			logger.Info("Saving Legacy Data to S3")
			s3.upload.uploads3("/tmp/" + string(m.Userid) + ".json")

			logger.Info("Posting to Watson")
			watson.WatsonPostCreateWallet(string(m.Userid))

			return profile, nil
		}
	}
}
