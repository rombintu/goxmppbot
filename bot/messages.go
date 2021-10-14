package bot

import (
	"errors"
	"fmt"
	"strings"
	"time"

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

// Action on /start
func OnStart() string {
	var submess string
	t := time.Now()
	switch {
	case t.Hour() < 12:
		submess = "Доброе утро"
	case t.Hour() < 17:
		submess = "Добрый день"
	default:
		submess = "Добрый вечер"
	}
	return submess + ", напишите название сервиса по которому у вас возник вопрос"
}

// Action on /help
func OnHelp() string {
	return "\nСервисы [с] - Вывести ссылки на ответы по сервисам\nПоддержка - Написать письмо в поддержку"
}

// Validate message for support mail
func ValidateSupport(message []string) error {
	// support:subject body
	size := len(message)
	// var subject, body string
	err := errors.New("Пример запроса должен быть типа: support:subject body")
	if size < 2 && message[0] != "support" || message[0] != "поддержка" {
		return err
	} else {
		inner := strings.Split(message[1], " ")
		if len(inner) < 2 {
			return err
		}
	}
	return nil
}

// Parse subject and message body from user text
func ParseSubjectAndBody(message []string) []string {
	// support:subject body
	var subject, body string
	inner := strings.Split(message[1], " ")
	subject = inner[0]
	body = strings.Join(inner[1:], " ")
	return []string{subject, body}
}

// Action on /support. Send mail to support
func (bot *Bot) OnSupport(subject, body string) (string, error) {
	if err := bot.SendToSupport(subject, body); err != nil {
		bot.Logger.Error(err)
		return "", err
	}
	return "Ваша заявка отправлена в обработку", nil
}

// Loop func, listening command from users
func (bot *Bot) HandleMessage() error {

	for {
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

			userText := strings.ToLower(data.(xmpp.Chat).Text)
			forSupport := strings.Split(userText, ":")
			if err := ValidateSupport(forSupport); err == nil {
				emailData := ParseSubjectAndBody(forSupport)
				resp, err := bot.OnSupport(emailData[0], emailData[1])
				if err != nil {
					mess.Text = "Произошла внутренняя ошибка: " + err.Error()
					continue
				}
				mess.Text = resp
				bot.SendMessage(mess)
				continue
			}
			switch userText {
			case "/start", "start", "старт":
				mess.Text = OnStart()
			case "/помощь", "помощь", "п", "help", "/help":
				mess.Text = OnHelp()
			case "list", "лист", "сервисы", "сервис", "с":
				buff := ""
				for key, value := range bot.Config.Links {
					buff = buff + fmt.Sprintf("%s -> %s\n", key, value)
				}
				mess.Text = buff
			case "поддержка":
				response, err := bot.OnSupport("", "")
				if err != nil {
					mess.Text = fmt.Sprint("Произошла ошибка: ", err.Error())
					bot.SendMessage(mess)
					continue
				}
				mess.Text = response
			default:
				mess.Text = "Неверный ввод, попробуйте ввести команду /help или помощь."
			}
			bot.SendMessage(mess)
		}
	}
}
