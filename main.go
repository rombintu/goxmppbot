package main

import (
	"os"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func main() {
	bot := xmppbot.NewBot()
	if err := bot.ConfigureBackand(); err != nil {
		bot.Logger.Error(err)
		os.Exit(1)
	}
	if err := bot.Connect(); err != nil {
		bot.Logger.Error(err)
		os.Exit(1)
	}
	if err := bot.HandleMessage(); err != nil {
		bot.Logger.Error(err)
		// os.Exit(1)
	}
}
