package createwallet_child_workflow

import (
	"context"
	"encoding/json"
	"os"

	kafka "github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

type Event2 struct {
	_airbyte_ab_id string `json:"-"`
	// json tag with - does not parse a result
	// items in struct not defined with a cap will not EXPORT or be available
	_airbyte_emitted_at string `json:"-"`
	_airbyte_data       struct {
		//Need to change topic name for each workflow
		Candidates_neocandidate struct {
			User_id string `json:"Userid"`
			//Need to change in env file best to setup go-env to handle at compile time
		}
	}
	_airbyte_stream string `json:"-"`
}

func CreateWalletParentWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "CREATEWALLET-WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string

	err := workflow.ExecuteChildWorkflow(ctx, CreateWalletChildWorkflow, "World").Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
	}

	logger.Info("Parent execution completed.", "Result", result)
	logger.Info("Writing message to next topic")
	b, _ := json.Marshal(result)
	e := Event2{}
	json.Unmarshal([]byte(b), &e)

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("KafkaURL")),
		Topic:        "Candidates_neocandidate",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 1,
	}

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("User_id"),
			Value: []byte(e._airbyte_data.Candidates_neocandidate.User_id),
		},
	)

	w.Close()
	// action to write to next topic
	return result, nil
}
