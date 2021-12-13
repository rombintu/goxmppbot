package main_test

import (
	"fmt"
	"testing"
	"time"

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
	if err := bot.SendToSupport("test@test.ru", "СУДИС", "тело текста -- тест test $#$FDF"); err != nil {
		t.Error(err)
	}
}
func TestSendToSupportManySymbols(t *testing.T) {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	if err := bot.SendToSupport("test@test.ru", "СУДИС", "gsdfffffffffffffffffffffffgdfsgdfgsdfgdsfgsdfgfdshfghfdggggggggggggggggggggggggggggggggggggggggggggggghgsfhsfghfghdfhfghfg"); err != nil {
		t.Error(err)
	}
}

func TestSendOOB(t *testing.T) {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	mess := xmppbot.CreateMessage()
	mess.Text = "test"
	mess.Remote = ""
	mess.Subject = "bothelper"
	mess.Ooburl = "http://risovach.ru/upload/2014/08/mem/nu-pochemu_59642579_orig_.jpg"
	mess.Stamp.Add(10 * time.Second)
	bot.SendMessage(mess)
}
