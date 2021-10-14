package main

import (
	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func main() {
	bot := xmppbot.NewBot()
	bot.Connect()
	bot.HandleMessage()
}
