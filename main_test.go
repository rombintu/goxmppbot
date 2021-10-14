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
