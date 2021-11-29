package bot

import (
	"errors"
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
		bot.Logger.Info(errorSend, err.Error())
	}

}

// Bot send message type-chat from chat-struct
func (bot *Bot) SendMessage(chat xmpp.Chat) {
	chat.Type = "chat"
	_, err := bot.Client.Send(chat)
	if err != nil {
		bot.Logger.Info(errorSend, err.Error())
	}
}

// Dev func
func (bot *Bot) SendOOB(chat xmpp.Chat) {
	_, err := bot.Client.SendOOB(chat)
	if err != nil {
		bot.Logger.Info(errorSend, err.Error())
	}
}

// Dev func
func (bot *Bot) SendORG(chat xmpp.Chat) {
	_, err := bot.Client.SendOrg("message")
	if err != nil {
		bot.Logger.Info(errorSend, err.Error())
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
		"help":       "помощь",
		"support":    "поддержка",
		"search":     "поиск",
		"start":      "старт",
		"services":   "сервисы",
		"refresh":    "/refresh",
		"addservice": "/addservice",
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
				bot.SendMessage(mess)
				continue
			}

			switch lastCommand {
			case reserv["search"]:
				mess.Text = loading
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
					bot.SendMessage(mess)
					continue
				}
				if len(resp) == 0 {
					mess.Text = notFound
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
					bot.SendMessage(mess)
					continue
				}
				mess.Text = resp
				bot.SendMessage(mess)
				continue
			case reserv["services"]:
				data, err := bot.Backend.GetJsonByName(ToLower(userText))
				if err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					bot.SendMessage(mess)
					continue
				}
				if len(data.Questions) == 0 {
					mess.Text = notFound
					bot.SendMessage(mess)
					continue
				}
				for i, q := range data.Questions {
					mess.Text = fmt.Sprintf("Вопрос: *%s*\n\tОтвет: %s\n ---\n", q.Subquestion[i], q.Subanswer[i])
					bot.SendMessage(mess)
				}
				continue
			case reserv["refresh"]:
				if userText != bot.Config.Default.RefreshSecret {
					continue
				}
				mess.Text = loading
				bot.SendMessage(mess)
				urls, _, err := bot.Backend.GetPageUrlsAndNames()
				if err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
				for _, u := range urls {
					page, err := GetPage(u)
					if err != nil {
						bot.Logger.Error(err)
						mess.Text = ToError(err)
						bot.SendMessage(mess)
						continue
					}
					if err := bot.Backend.UpdatePage(page, u); err != nil {
						bot.Logger.Error(err)
						mess.Text = ToError(err)
						bot.SendMessage(mess)
						continue
					}
				}
				mess.Text = dbUpdated
				bot.SendMessage(mess)
				continue
			case reserv["addservice"]:
				text := strings.Split(userText, "|")
				if text[0] != bot.Config.Default.RefreshSecret {
					continue
				}
				mess.Text = loading
				bot.SendMessage(mess)
				if len(text) != 3 {
					mess.Text = ToError(errors.New(fewArguments))
					continue
				}
				if err := bot.Backend.PutNewPage(text[1], text[2]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
				mess.Text = dbUpdated
				bot.SendMessage(mess)
				continue
			}

			switch ToLower(userText) {
			case "/start", reserv["start"]:
				mess.Text = OnStart()
			case "/помощь", "/help", reserv["help"]:
				mess.Text = onHelp
			case reserv["services"]:
				buff := serviceHelpMessage
				_, names, err := bot.Backend.GetPageUrlsAndNames()
				if err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}

				for i, n := range names {
					buff += fmt.Sprintf("\t%d. %s\n", i+1, strings.ToTitle(n))
				}
				mess.Text = buff
				if err := bot.Backend.PutCommand(GetHash(from), reserv["services"]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
			case reserv["support"]:
				mess.Text = supportHelpMessage
				if err := bot.Backend.PutCommand(GetHash(from), reserv["support"]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
			case reserv["search"]:
				mess.Text = searchHelpMessage
				if err := bot.Backend.PutCommand(GetHash(from), reserv["search"]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
			case reserv["refresh"]:
				if err := bot.Backend.PutCommand(GetHash(from), reserv["refresh"]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
				mess.Text = "Enter password"
			case reserv["addservice"]:
				if err := bot.Backend.PutCommand(GetHash(from), reserv["addservice"]); err != nil {
					bot.Logger.Error(err)
					mess.Text = ToError(err)
					continue
				}
				mess.Text = "password|name|url"
			case "last":
				if lastCommand == "" {
					mess.Text = "Null"
				} else {
					mess.Text = lastCommand
				}
			case "":
				continue
			default:
				mess.Text = notFoundCommand
			}
			bot.SendMessage(mess)
		}
		if err := bot.HandleMessage(); err != nil {
			return err
		}
	}
}
