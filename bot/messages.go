package bot

import (
	"log"

	xmpp "github.com/mattn/go-xmpp"
)

func CreateMessage(to, subject, message string) xmpp.Chat {
	return xmpp.Chat{
		Remote:  to,
		Text:    message,
		Subject: subject,
	}
}

func (bot *Bot) SendEmail(chat xmpp.Chat) {
	_, err := bot.Client.Send(chat)
	if err != nil {
		log.Println("Error send message: ", err.Error())
	}

}

func (bot *Bot) SendMessage(chat xmpp.Chat) {
	chat.Type = "chat"
	_, err := bot.Client.Send(chat)
	if err != nil {
		log.Println("Error send message: ", err.Error())
	}
}

func (bot *Bot) HandleMessage() error {

	for {
		//
		data, err := bot.Client.Recv()
		if err != nil {
			return err
		}
		switch data.(type) {
		case xmpp.Chat:
			log.Printf("%+v", data)

		}

	}
}
