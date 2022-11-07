package child_workflow

import (
	"go.temporal.io/sdk/workflow"
)

// Kafka listen to topic
// get message from topic
// parse message for user_id
// POST to localhost:7001/wallet/user_id

func SampleChildWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	greeting := "Hello " + name + "!"
	logger.Info("Child workflow execution: " + greeting)
	return greeting, nil
}

// @@@SNIPEND
