package bot

import (
	"log"

	xmpp "github.com/mattn/go-xmpp"
)

type Bot struct {
	Host     string
	Login    string
	Password string
	Debug    bool
	Client   *xmpp.Client
}

func (bot *Bot) Connect() {
	client, err := xmpp.NewClientNoTLS(bot.Host, bot.Login, bot.Password, true)
	if err != nil {
		log.Fatal("Error connect: ", err.Error())
	}
	bot.Client = client
}
