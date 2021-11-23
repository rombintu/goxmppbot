package main

import (
	"os"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func main() {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		bot.Logger.Error(err)
		os.Exit(1)
	}
	if err := bot.HandleMessage(); err != nil {
		bot.Logger.Error(err)
		os.Exit(1)
	}
}
