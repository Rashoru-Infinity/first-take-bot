package usecase

import (
	"github.com/Rashoru-Infinity/first-take-bot/domain"
)

type Event struct {
}

type EventIF interface {
	RegisterEvent(handler *domain.SocketmodeHandler, ec domain.EventDriverIF)
	RunEventLoop(handler *domain.SocketmodeHandler)
}

func NewEvent() EventIF {
	return Event{}
}

func (e Event) RegisterEvent(handler *domain.SocketmodeHandler, ec domain.EventDriverIF) {
	ec.RegisterMessageEventHandler(handler, ec)
	ec.RegisterConnectingEventHandler(handler, ec)
	ec.RegisterConnectionErrorEventHandler(handler, ec)
	ec.RegisterConnectedEventHandler(handler, ec)
}

func (e Event) RunEventLoop(handler *domain.SocketmodeHandler) {
	handler.Handler.RunEventLoop()
}
