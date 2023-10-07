package domain

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type EventDriver struct {
}

type EventDriverIF interface {
	GetEvent(evt *socketmode.Event) (slackevents.EventsAPIEvent, bool)
	GetMessageEvent(eventsAPIEvent slackevents.EventsAPIEvent) (*slackevents.MessageEvent, bool)
	PostMessage(channel, msg string, escape bool, sockclient *socketmode.Client) error
	HandleMessageEvent(evt *socketmode.Event, sockclient *socketmode.Client)
	MiddlewareConnecting(evt *socketmode.Event, client *socketmode.Client)
	MiddlewareConnectionError(evt *socketmode.Event, client *socketmode.Client)
	MiddlewareConnected(evt *socketmode.Event, client *socketmode.Client)
	RegisterMessageEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF)
	RegisterConnectingEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF)
	RegisterConnectionErrorEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF)
	RegisterConnectedEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF)
}

func NewEventDriver() EventDriverIF {
	return EventDriver{}
}

func (ed EventDriver) GetEvent(evt *socketmode.Event) (slackevents.EventsAPIEvent, bool) {
	eventsAPIEVENT, ok := evt.Data.(slackevents.EventsAPIEvent)
	return eventsAPIEVENT, ok
}

func (ed EventDriver) GetMessageEvent(eventsAPIEvent slackevents.EventsAPIEvent) (*slackevents.MessageEvent, bool) {
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	return ev, ok
}

func (ed EventDriver) PostMessage(channel, msg string, escape bool, sockclient *socketmode.Client) error {
	_, _, err := sockclient.Client.PostMessage(channel, slack.MsgOptionText(msg, escape))
	return err
}

func (ed EventDriver) HandleMessageEvent(evt *socketmode.Event, sockclient *socketmode.Client) {
	var (
		payload map[string]interface{}
		author  string
		content string
		channel = os.Getenv("BOT_CHANNEL")
		selfID  = os.Getenv("BOT_SELFID")
		msgesc  bool
	)
	msgesc, err := strconv.ParseBool(os.Getenv("BOT_MESSAGE_ESCAPE"))
	if err != nil {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	eventsAPIEvent, ok := ed.GetEvent(evt)
	if !ok {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	sockclient.Ack(*evt.Request)
	ev, ok := ed.GetMessageEvent(eventsAPIEvent)
	if !ok {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	if err := json.Unmarshal(evt.Request.Payload, &payload); err != nil {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	author, ok = payload["event"].(map[string]interface{})["user"].(string)
	if !ok {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	if author == selfID {
		log.Println("skipped self message")
		return
	}
	content, ok = payload["event"].(map[string]interface{})["text"].(string)
	if !ok {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	err = ed.PostMessage(channel, author+"@"+ev.Channel+": "+content, msgesc, sockclient)
	if err != nil {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	files, ok := payload["event"].(map[string]interface{})["files"].([]interface{})
	if !ok {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
	urls := make([]byte, 0, 10)
	for _, f := range files {
		url, ok := f.(map[string]interface{})["url_private"].(string)
		if !ok {
			pc, filename, line, _ := runtime.Caller(0)
			log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
			return
		}
		urls = append(urls, url...)
		urls = append(urls, 0x20)
	}
	err = ed.PostMessage(channel, author+"@"+ev.Channel+": "+string(urls), false, sockclient)
	if err != nil {
		pc, filename, line, _ := runtime.Caller(0)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		return
	}
}

func (ed EventDriver) MiddlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("connecting to slack with socket mode...")
}

func (ed EventDriver) MiddlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	log.Println(evt.Data)
	log.Println("connection failed. retrying later...")
}

func (ed EventDriver) MiddlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("connected to slack with socket mode.")
}

func (ed EventDriver) RegisterMessageEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF) {
	handler.Handler.HandleEvents(slackevents.Message, handlerFunc.HandleMessageEvent)
}

func (ed EventDriver) RegisterConnectingEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF) {
	handler.Handler.Handle(socketmode.EventTypeConnecting, handlerFunc.MiddlewareConnecting)
}

func (ed EventDriver) RegisterConnectionErrorEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF) {
	handler.Handler.Handle(socketmode.EventTypeConnectionError, handlerFunc.MiddlewareConnectionError)
}

func (ed EventDriver) RegisterConnectedEventHandler(handler *SocketmodeHandler, handlerFunc EventDriverIF) {
	handler.Handler.Handle(socketmode.EventTypeConnected, handlerFunc.MiddlewareConnected)
}
