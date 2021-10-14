package main

import (
	"os"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func main() {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		os.Exit(1)
	}
	bot.HandleMessage()
}
