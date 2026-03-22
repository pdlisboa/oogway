package client

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
)

type BedrockClient struct {
	client *bedrockagentruntime.Client
}

func NewBedrockCLient() (AgentClient, error) {
	var bedrockagentRuntimeClient *bedrockagentruntime.Client
	sync.OnceFunc(func() {
		ctx := context.TODO()
		sdkConfig, initError := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
		if initError != nil {
			slog.Error("Error loading default config", "error", initError)
			return
		}

		bedrockagentRuntimeClient = bedrockagentruntime.NewFromConfig(sdkConfig)
	})
	return BedrockClient{
		client: bedrockagentRuntimeClient,
	}, nil
}

func (br BedrockClient) Send(ctx context.Context, req string) (string, error) {

	if br.client == nil {
		return "", errors.New("Client not initialized")
	}

	params := &bedrockagentruntime.InvokeAgentInput{
		AgentId:      aws.String("YourAgentID"),
		AgentAliasId: aws.String("YourAgentAliasID"),
		InputText:    aws.String(req),
		SessionId:    aws.String("RandomSessionID"),
	}
	invoke, err := br.client.InvokeAgent(context.Background(), params)

	if err != nil {
		return "", err
	}

	var result string
	for event := range invoke.GetStream().Events() {
		switch e := event.(type) {
		case *types.ResponseStreamMemberChunk:
			result = string(e.Value.Bytes)
		}
	}

	return result, nil
}
