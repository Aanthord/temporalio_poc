package createwallet_child_workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aanthord/temporalio_poc/watson"
	kafka "github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	service     = "temporalio-createwallet"
	environment = "test"
	id          = 1
)

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

func CreateWalletChildWorkflow(ctx workflow.Context, name string) (string, error) {
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

	w.RegisterWorkflow(createwallet_child_workflow.CreateWalletParentWorkflow)
	w.RegisterWorkflow(createwallet_child_workflow.CreateWalletChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
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
		b, _ := json.Marshal(m)
		json.Unmarshal([]byte(b), &payload)    // Convert JSON data into interface{} type
		um := payload.(map[string]interface{}) // To use the converted data we will need to convert it
		// into a map[string]interface{}

		logger.Info("Getting user_id")

		//Need to do stuff here so I can pass userID to watson
		logger.Info("Posting to Watson")
		watson.WatsonPostCreateWallet(um.Userid)

	}
}
