package client

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/openrouter"
)

type OpenrouterClient struct {
	client fantasy.Agent
}

func NewOpenrouterCLient(apiKey string, headers map[string]string) (AgentClient, error) {
	opts := []openrouter.Option{
		openrouter.WithAPIKey(apiKey),
	}

	if len(headers) > 0 {
		opts = append(opts, openrouter.WithHeaders(headers))
	}

	provider, err := openrouter.New(opts...)

	if err != nil {
		slog.Error("Error on building openrouter provider connection", "error", err)
		return nil, err
	}

	ctx := context.Background()

	model, err := provider.LanguageModel(ctx, "minimax/minimax-m2.5:free")
	if err != nil {
		slog.Error("Error on geting openrouter LM ", "error", err)
		os.Exit(1)
	}

	agent := fantasy.NewAgent(
		model,
		fantasy.WithSystemPrompt("You are a senior programmer specialist in Golang. Your name is oogway!"),
	)

	return OpenrouterClient{
		client: agent,
	}, nil
}

func (or OpenrouterClient) Send(ctx context.Context, req string) (string, error) {

	if or.client == nil {
		return "", errors.New("Client not initialized")
	}

	result, err := or.client.Generate(ctx, fantasy.AgentCall{Prompt: req})

	if err != nil {
		slog.Error("Erro on response ", "error", err)
		os.Exit(1)
	}

	return result.Response.Content.Text(), nil
}
