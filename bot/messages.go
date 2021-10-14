package bot

import (
	"fmt"
	"strings"

	xmpp "github.com/mattn/go-xmpp"
)

// Return message chat-struct
func CreateMessage() xmpp.Chat {
	return xmpp.Chat{}
}

// Bot send message type-email from chat-struct
func (bot *Bot) SendEmail(chat xmpp.Chat) {
	chat.Type = "normal"
	_, err := bot.Client.Send(chat)
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}

}

// Bot send message type-chat from chat-struct
func (bot *Bot) SendMessage(chat xmpp.Chat) {
	chat.Type = "chat"
	_, err := bot.Client.Send(chat)
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}
}

// Dev func
func (bot *Bot) SendOOB(chat xmpp.Chat) {
	_, err := bot.Client.SendOOB(chat)
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}
}

// Dev func
func (bot *Bot) SendORG(chat xmpp.Chat) {
	_, err := bot.Client.SendOrg("message")
	if err != nil {
		bot.Logger.Info("Error send message: ", err.Error())
	}
}

// Loop func, listening command from users
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
			case "/start", "start", "старт":
				mess.Text = "Привет, напишите название сервиса по которому у вас возник вопрос"
			case "/помощь", "помощь", "п", "help", "/help":
				mess.Text = "\nСервисы [с] - Вывести ссылки на ответы по сервисам\nПоддержка - Написать письмо в поддержку"
			// case "сэп":
			// 	mess.Text = "http://it.mvd.ru/sections/4#questions"
			// case "календарь":
			// 	mess.Text = "Выберите из списка:\n1. ...\n2. ... \n3. ...\n4. Прочее"
			case "list", "лист", "сервисы", "сервис", "с":
				buff := ""
				for key, value := range bot.Config.Links {
					buff = buff + fmt.Sprintf("%s -> %s\n", key, value)
				}
				mess.Text = buff
			case "поддержка":
				subject := "botThemeTest"
				bodyMessage := "bodyBotTest"
				err := bot.SendToSupport(subject, bodyMessage)
				if err != nil {
					bot.Logger.Error(err)
				}
				mess.Text = "Отправлено!"
			default:
				mess.Text = "Неверный ввод, попробуйте ввести команду /help или помощь."
			}
			bot.SendMessage(mess)
		}
	}
}
