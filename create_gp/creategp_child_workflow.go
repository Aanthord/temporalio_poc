package child_workflow

import (
	"go.temporal.io/sdk/workflow"
)

func CreateRecruiterChildWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	greeting := "Hello " + name + "!"
	logger.Info("Child workflow execution: " + greeting)
	return greeting, nil
}
