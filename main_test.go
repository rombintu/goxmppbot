package main_test

import (
	"fmt"
	"testing"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func TestParseConfig(t *testing.T) {
	bot := xmppbot.NewBot()

	fmt.Println(bot.Config.Links)
}

func TestSendToSupport(t *testing.T) {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	if err := bot.SendToSupport("СУДИС", "тело текста -- тест test $#$FDF"); err != nil {
		t.Error(err)
	}
}
