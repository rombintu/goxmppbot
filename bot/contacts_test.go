package bot_test

import (
	"fmt"
	"testing"

	"github.com/rombintu/goxmppbot/bot"
)

func TestGetUserByID(t *testing.T) {
	data, err := bot.GetUserByRegex("ashetukhin*", "https://mail.leaguemail.ru/contacts", "10")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}
