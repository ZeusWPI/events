package zauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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

	html, err := toHTML(content)
	if err != nil {
		return err
	}

	m := mail{
		Subject: subject,
		Content: html,
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

func toHTML(content string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", fmt.Errorf("convert markdown to html %s | %w", content, err)
	}

	return buf.String(), nil
}
