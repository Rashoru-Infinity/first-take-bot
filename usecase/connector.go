package usecase

import (
	"errors"

	"github.com/Rashoru-Infinity/first-take-bot/domain"
)

type Connector struct {
}

type ConnectorIF interface {
	ConnectSlack(apptoken, bottoken string, sd domain.SlackDriverIF) (*domain.SocketmodeHandler, error)
}

func NewConnector() ConnectorIF {
	return Connector{}
}

func (c Connector) ConnectSlack(apptoken, bottoken string, sd domain.SlackDriverIF) (*domain.SocketmodeHandler, error) {
	appLevelToken := sd.GetOptionAppLevelToken(apptoken)
	if appLevelToken == nil {
		return nil, errors.New("GetOptionAppLevelToken")
	}
	api := sd.NewAPI(bottoken, appLevelToken)
	if api == nil {
		return nil, errors.New("NewAPI")
	}
	client := sd.NewClient(api)
	if client == nil {
		return nil, errors.New("NewClient")
	}
	handler := sd.NewSocketmodeHandler(client)
	if handler == nil {
		return nil, errors.New("NewSocketmodeHandler")
	}
	return handler, nil
}
