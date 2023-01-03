package createorgwallet_child_workflow

import (
	"go.temporal.io/sdk/workflow"
)

func CreateOrgWalletParentWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "CREATEORGWALLET-WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string
	err := workflow.ExecuteChildWorkflow(ctx, CreateOrgWalletChildWorkflow, "World").Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return "", err
	}

	logger.Info("Parent execution completed.", "Result", result)
	logger.Info("Writing message to next topic")
	//kfka.Writer(validateWallet, "userid", m.["userid"].(string))
	// action to write to next topic
	return result, nil
}
