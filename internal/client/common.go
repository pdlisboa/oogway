package client

import "context"

type AgentClient interface {
	Send(ctx context.Context, req string) (string, error)
}

type ResponseMsg struct {
	Content string
	Err     error
}
