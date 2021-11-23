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

// Action on /support. Send mail to support
func (bot *Bot) OnSupport(user, subject, body string) (string, error) {
	if err := bot.SendToSupport(user, subject, body); err != nil {
		bot.Logger.Error(err)
		return "", err
	}
	return "Ваша заявка отправлена в обработку", nil
}

// Loop func, listening command from users
func (bot *Bot) HandleMessage() error {
	reserv := map[string]string{
		"help":     "помощь",
		"support":  "поддержка",
		"search":   "поиск",
		"start":    "старт",
		"services": "сервисы",
	}

	for {
		data, err := bot.Client.Recv()
		if err != nil {
			return err
		}
		bot.Logger.Info(data)
		switch data.(type) {
		case xmpp.Chat:
			mess := CreateMessage()
			from := data.(xmpp.Chat).Remote

			mess.Remote = from
			// Theme for Chat (not email)
			mess.Subject = "bothelper"

			userText := data.(xmpp.Chat).Text

			if userText == "" || userText == " " {
				continue
			}

			lastCommand, err := bot.Backend.GetLastCommand(GetHash(from))
			if err != nil {
				bot.Logger.Error(err)
				mess.Text = ToError(err)
				continue
			}

			switch lastCommand {
			case reserv["search"]:
				mess.Text = "Выполняется..."
				bot.SendMessage(mess)
				count := "5"
				target := strings.Split(ToLower(userText), ":")
				if len(target) > 1 {
					count = TrimS(target[1])
				}
				resp, err := GetUserByRegex(target[0], bot.Config.Contacts.Url, count)
				if err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
				if len(resp) == 0 {
					mess.Text = "Ничего не найдено"
					bot.SendMessage(mess)
					continue
				}
				mess.Text = BuildMessageFromUsers(resp)
				bot.SendMessage(mess)
				continue
			case reserv["support"]:
				emailData, err := ParseSubjectAndBody(userText)
				if err != nil {
					mess.Text = err.Error()
					bot.SendMessage(mess)
					continue
				}
				resp, err := bot.OnSupport(from, emailData[0], emailData[1])
				if err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
				mess.Text = resp
				bot.SendMessage(mess)
				continue
			}

			switch ToLower(userText) {
			case "/start", reserv["start"]:
				mess.Text = OnStart()
			case "/помощь", "/help", reserv["help"]:
				mess.Text = OnHelp()
			case "list", "лист", reserv["services"]:
				buff := ""
				for key, value := range bot.Config.Links {
					buff = buff + fmt.Sprintf("%s [%s]\n", value, key)
				}
				mess.Text = buff
			case reserv["support"]:
				mess.Text = `
Напишите свое обращение такого вида:
Поддержка:НАЗВАНИЕ_СЕРВИСА письмо
	Пример: *СУДИС Все сломалось, помогите*
			`
				if err := bot.Backend.PutCommand(GetHash(from), reserv["support"]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
			case reserv["search"]:
				mess.Text = `
Напишите свое обращение такого вида:
Поиск: ФИО_ПОЧТА_ДОЛЖНОСТЬ_КОМПАНИЯ
Примечание: Можно использовать регулярные выражения
Примечание: Добавьте в конце *: N*, чтобы регулировать выборку
	1. Пример: *Иванов*
	2. Пример: *Иванов: 10*
	3. Пример: *ivanov*
			`
				if err := bot.Backend.PutCommand(GetHash(from), reserv["search"]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
			case "last":
				c, err := bot.Backend.GetLastCommand(GetHash(from))
				if err != nil {
					return err
				}
				if c == "" {
					mess.Text = "Null"
				} else {
					mess.Text = c
				}
			case "":
				continue
			default:
				mess.Text = "Неверный ввод, попробуйте ввести команду [help] или [помощь]."
			}
			bot.SendMessage(mess)
		}
		if err := bot.HandleMessage(); err != nil {
			return err
		}
	}
}
