package lagout

import (
	"context"
)

// Client is a chat client stub (e.g. Lagout or similar). Per mandated path: chat/lagout/client.go.
type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Send(ctx context.Context, channel, message string) error {
	_ = ctx
	_, _ = channel, message
	return nil
}
