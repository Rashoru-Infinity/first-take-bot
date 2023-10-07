package main

import (
	"log"
	"os"

	"github.com/Rashoru-Infinity/first-take-bot/domain"
	"github.com/Rashoru-Infinity/first-take-bot/usecase"
)

func main() {
	var (
		slackAppToken = os.Getenv("APP_TOKEN")
		slackBotToken = os.Getenv("BOT_TOKEN")
	)
	slackDriver := domain.NewSlackDriver()
	connector := usecase.NewConnector()
	sockmodeHandler, err := connector.ConnectSlack(slackAppToken, slackBotToken, slackDriver)
	if err != nil {
		log.Println(err)
		return
	}
	eventDriver := domain.NewEventDriver()
	event := usecase.NewEvent()
	event.RegisterEvent(sockmodeHandler, eventDriver)
	event.RunEventLoop(sockmodeHandler)
}
