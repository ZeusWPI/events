package mattermost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type Message struct {
	ChannelID string `json:"channel_id"`
	Message   string `json:"message"`
}

type messageResponse struct {
	// Don't care about the response at this time
}

func (c *Client) SendMessage(ctx context.Context, message Message) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(message)
	if err != nil {
		return fmt.Errorf("encode message %+v | %w", message, err)
	}

	query := query{
		method: "POST",
		url:    "posts",
		body:   &buf,
		target: new(messageResponse),
	}

	if err := c.query(ctx, query); err != nil {
		return err
	}

	return nil
}
