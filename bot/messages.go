package bot

import (
	"strings"

	xmpp "github.com/mattn/go-xmpp"
)

func CreateMessage() xmpp.Chat {
	return xmpp.Chat{}
}

func (bot *Bot) SendEmail(chat xmpp.Chat) {
	chat.Type = "normal"
	_, err := bot.Client.Send(chat)
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}

}

func (bot *Bot) SendMessage(chat xmpp.Chat) {
	chat.Type = "chat"
	_, err := bot.Client.Send(chat)
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}
}

func (bot *Bot) SendOOB(chat xmpp.Chat) {
	// chat.Type = "normal"
	_, err := bot.Client.SendOOB(chat)
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}
}

func (bot *Bot) SendORG(chat xmpp.Chat) {
	// chat.Type = "normal"
	_, err := bot.Client.SendOrg("message")
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}
}

func (bot *Bot) HandleMessage() error {

	for {
		//
		data, err := bot.Client.Recv()
		if err != nil {
			return err
		}
		bot.Logger.Info(data)
		switch data.(type) {
		case xmpp.Chat:
			mess := CreateMessage()
			mess.Remote = data.(xmpp.Chat).Remote
			mess.Subject = "bothelper"
			switch strings.ToLower(data.(xmpp.Chat).Text) {
			case "/start":
				mess.Text = "Привет, напишите название сервиса по которому у вас возник вопрос"
			case "сэп":
				mess.Text = "Какой раздел?"
			case "календарь":
				mess.Text = "Выберите из списка:\n1. ...\n2. ... \n3. ...\n4. Прочее"
			case "прочее":
				mess.Text = "Отправлено!"

			default:
				mess.Text = "Не полнял ввод, попробуйте снова."
			}
			bot.SendMessage(mess)
		}
	}
}
