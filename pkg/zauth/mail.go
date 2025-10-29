package zauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type mail struct {
	Subject string `json:"subject"`
	Content string `json:"body"`
	Author  string `json:"author"`
}

func (c *client) MailAll(ctx context.Context, subject, content string) error {
	conf := &clientcredentials.Config{
		ClientID:     c.clientKey,
		ClientSecret: c.secret,
		TokenURL:     endpoint + "/oauth/token",
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	client := conf.Client(ctx)

	m := mail{
		Subject: subject,
		Content: content,
		Author:  "Zeus WPI",
	}
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal mail bytes %+v | %w", m, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/mails", bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http client post %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status http code %+v | %s", *resp, resp.Status)
	}

	return nil
}
