package creategp_child_workflow

import (
	"go.temporal.io/sdk/workflow"
)

func CreateGPParentWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "CREATEGP-WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string
	err := workflow.ExecuteChildWorkflow(ctx, CreateGPChildWorkflow, "World").Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return "", err
	}

	logger.Info("Parent execution completed.", "Result", result)
	logger.Info("Writing message to next topic")
	//kfka.Writer(validateWallet, "userid", string(m.Userid))
	return result, nil
}
