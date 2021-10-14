package main_test

import (
	"fmt"
	"testing"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func TestSendToSupport(t *testing.T) {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	if err := bot.SendToSupport("themeTest", "bodyTest"); err != nil {
		t.Error(err)
	}
}

func TestParseConfig(t *testing.T) {
	bot := xmppbot.NewBot()

	fmt.Println(bot.Config.Links)
}
