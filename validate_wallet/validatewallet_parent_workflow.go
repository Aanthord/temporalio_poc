package validatewallet_child_workflow

import (
	"context"
	"encoding/json"
	"os"

	kafka "github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

type Event1 struct {
	User_id string `json:"Userid"`
}

func ValidateWalletParentWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "VALIDATEWALLET-WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string
	err := workflow.ExecuteChildWorkflow(ctx, ValidateWalletChildWorkflow, "World").Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return "", err
	}

	logger.Info("Parent execution completed.", "Result", result)
	b, _ := json.Marshal(result)
	e := Event1{}
	json.Unmarshal([]byte(b), &e)

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("KafkaURL")),
		Topic:        "CreateCandidateNFT",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 1,
	}

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("User_id"),
			Value: []byte(e.User_id),
		},
	)

	w.Close()
	// action to write to next topic
	return result, nil
}
