package zauth

import "github.com/ZeusWPI/events/pkg/config"

const endpoint = "https://zauth.zeus.gent"

var C *client

type client struct {
	clientKey   string
	secret      string
	callbackURL string
}

func Init() {
	C = &client{
		clientKey:   config.GetString("auth.client"),
		secret:      config.GetString("auth.secret"),
		callbackURL: config.GetString("auth.callback_url"),
	}
}
