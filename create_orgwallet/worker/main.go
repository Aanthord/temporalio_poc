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

	child_workflow "github.com/aanthord/temporalio_poc/create_orgwallet"
)
//get a list of all the topics
func getKafkaTopics(kafkaURL) {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
    panic(err.Error())
}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
    		panic(err.Error())
}

	m := map[string]struct{}{}

	for _, p := range partitions {
    	m[p.Topic] = struct{}{}
}
	for k := range m {
    	fmt.Println(k)
}
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
//Brain hurt need to iterate over topics in the list and call main for each topic creating a distinct
// count of messages in each topic. 

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

	w.RegisterWorkflow(child_workflow.CreateWalletParentWorkflow)
	w.RegisterWorkflow(child_workflow.CreateWalletChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
			log.Fatalln("Unable to start worker", err)
	}
}

}
