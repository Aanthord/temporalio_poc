package main

import (
	"log"
	"context"
	"format"
	"os"
	"stings"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	kafka "github.com/segmentio/kafka-go"

	child_workflow "github.com/aanthord/temporalio_poc/child-workflow"
)

func main() {
	// The client is a heavyweight object that should be created only once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "child-workflow", worker.Options{})

	w.RegisterWorkflow(child_workflow.SampleParentWorkflow)
	w.RegisterWorkflow(child_workflow.SampleChildWorkflow)

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
		
	}
}
