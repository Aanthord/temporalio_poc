package child_workflow

import (
	"go.temporal.io/sdk/workflow"
)

// CREATE_WALLET
// Listen to kafka topic
// Get message from topic
// parse message for user_id
// POST user_id to localhost:7001/wallet

func CreateWalletChildWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	greeting := "Hello " + name + "!"
	logger.Info("Child workflow execution: " + greeting)
	return greeting, nil
}

// @@@SNIPEND
