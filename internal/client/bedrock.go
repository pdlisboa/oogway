package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/google/uuid"
)

const (
	defaultBedrockRegion = "us-east-1"
	defaultAgentID       = "xxxx"
	defaultAgentAliasID  = "xxxxx"
)

type BedrockClient struct {
	client       *bedrockagentruntime.Client
	agentID      string
	agentAliasID string
	sessionID    string
}

func NewBedrockClient() (AgentClient, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(defaultBedrockRegion))
	if err != nil {
		slog.Error("Error loading default AWS config", "error", err)
		return nil, fmt.Errorf("load default AWS config: %w", err)
	}

	return BedrockClient{
		client:       bedrockagentruntime.NewFromConfig(sdkConfig),
		agentID:      defaultAgentID,
		agentAliasID: defaultAgentAliasID,
		sessionID:    uuid.NewString(),
	}, nil
}

func (br BedrockClient) Send(ctx context.Context, req string) (string, error) {
	if br.client == nil {
		return "", errors.New("client not initialized")
	}

	if strings.TrimSpace(req) == "" {
		return "", errors.New("request cannot be empty")
	}

	params := &bedrockagentruntime.InvokeAgentInput{
		AgentId:      aws.String(br.agentID),
		AgentAliasId: aws.String(br.agentAliasID),
		InputText:    aws.String(req),
		SessionId:    aws.String(br.sessionID),
	}

	slog.Debug("Bedrock InvokeAgent request",
		"agentID", br.agentID,
		"aliasID", br.agentAliasID,
		"sessionID", br.sessionID,
		"inputLen", len(req))

	invoke, err := br.client.InvokeAgent(ctx, params)

	if err != nil {
		return "", err
	}
	defer invoke.GetStream().Close()

	var result strings.Builder
	for event := range invoke.GetStream().Events() {
		switch e := event.(type) {
		case *types.ResponseStreamMemberChunk:
			_, err = result.Write(e.Value.Bytes)
			if err != nil {
				return "", fmt.Errorf("write response chunk: %w", err)
			}
		case *types.UnknownUnionMember:
			slog.Warn("Unknown Bedrock event received", "tag", e.Tag)
		}
	}

	if err := invoke.GetStream().Err(); err != nil {
		return "", fmt.Errorf("read response stream: %w", err)
	}

	return result.String(), nil
}
