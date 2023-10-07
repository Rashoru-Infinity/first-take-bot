package domain

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type AppLevelToken struct {
	Token func(*slack.Client)
}

type API struct {
	API *slack.Client
}

type Client struct {
	Client *socketmode.Client
}

type SocketmodeHandler struct {
	Handler *socketmode.SocketmodeHandler
}

type SlackDriver struct {
}

type SlackDriverIF interface {
	GetOptionAppLevelToken(bottoken string) *AppLevelToken
	NewAPI(appToken string, appLevelToken *AppLevelToken) *API
	NewClient(api *API) *Client
	NewSocketmodeHandler(client *Client) *SocketmodeHandler
}

func NewSlackDriver() SlackDriverIF {
	return SlackDriver{}
}

func (sc SlackDriver) GetOptionAppLevelToken(bottoken string) *AppLevelToken {
	token := slack.OptionAppLevelToken(bottoken)
	if token == nil {
		return nil
	}
	return &AppLevelToken{
		Token: token,
	}
}

func (sc SlackDriver) NewAPI(appToken string, appLevelToken *AppLevelToken) *API {
	api := slack.New(
		appToken,
		appLevelToken.Token,
	)
	if api == nil {
		return nil
	}
	return &API{
		API: api,
	}
}

func (sc SlackDriver) NewClient(api *API) *Client {
	client := socketmode.New(
		api.API,
	)
	if client == nil {
		return nil
	}
	return &Client{
		Client: client,
	}
}

func (sc SlackDriver) NewSocketmodeHandler(client *Client) *SocketmodeHandler {
	handler := socketmode.NewSocketmodeHandler(client.Client)
	if handler == nil {
		return nil
	}
	return &SocketmodeHandler{
		Handler: handler,
	}
}
