package bot

import (
	"encoding/json"
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

func (bot *Bot) SendError(err error) {
	mess := CreateMessage()
	mess.Text = err.Error()
	bot.SendMessage(mess)
}

// Loop func, listening command from users
func (bot *Bot) HandleMessage() {

	for {
		data, err := bot.Client.Recv()
		if err != nil {
			bot.Logger.Error(err)
		}

		switch data.(type) {
		case xmpp.Chat:
			if data.(xmpp.Chat).Text == "" || data.(xmpp.Chat).Text == " " {
				continue
			}
			if err := bot.Run(data); err != nil {
				bot.Logger.Error(err)
			}
		}
	}
}

func (bot *Bot) Run(data interface{}) error {
	bot.Logger.Info(data)
	mess := CreateMessage()
	from := data.(xmpp.Chat).Remote
	mess.Remote = from
	mess.Subject = "bothelper"
	userText := data.(xmpp.Chat).Text

	lastCommand, err := bot.Backend.GetLastCommand(GetHash(from))
	if err != nil {
		bot.SendError(err)
		return err
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
			bot.SendError(err)
			return err
		}
		if len(resp) == 0 {
			mess.Text = notFound
			bot.SendMessage(mess)
			return err
		}
		mess.Text = BuildMessageFromUsers(resp)
		bot.SendMessage(mess)
	case reserv["support"]:
		emailData, err := ParseSubjectAndBody(userText)
		if err != nil {
			bot.SendError(err)
			return err
		}
		resp, err := bot.OnSupport(from, emailData[0], emailData[1])
		if err != nil {
			bot.SendError(err)
			return err
		}
		mess.Text = resp
		bot.SendMessage(mess)
	case reserv["services"]:
		data, err := bot.Backend.GetJsonByName(ToLower(userText))
		if err != nil {
			bot.SendError(err)
			return err
		}

		if len(data) == 0 {
			mess.Text = notFound
			bot.SendMessage(mess)
			return nil
		}
		var page Page
		if err := json.Unmarshal(data, &page); err != nil {
			bot.SendError(err)
			return err
		}

		for i, q := range page.Questions {
			mess.Text = fmt.Sprintf("Вопрос: *%s*\n\tОтвет: %s\n ---\n", q.Subquestion[i], q.Subanswer[i])
			bot.SendMessage(mess)
		}
	case reserv["refresh"]:
		if userText != bot.Config.Default.RefreshSecret {
			return nil
		}
		mess.Text = loading
		bot.SendMessage(mess)
		urls, _, err := bot.Backend.GetPageUrlsAndNames()
		if err != nil {
			bot.SendError(err)
			return err
		}
		for _, u := range urls {
			page, err := GetPage(u)
			if err != nil {
				bot.SendError(err)
				return err
			}
			if err := bot.Backend.UpdatePage(page, u); err != nil {
				bot.SendError(err)
				return err
			}
		}
		mess.Text = dbUpdated
		bot.SendMessage(mess)
	case reserv["addservice"]:
		text := strings.Split(userText, "|")
		if text[0] != bot.Config.Default.RefreshSecret {
			return nil
		}
		mess.Text = loading
		bot.SendMessage(mess)
		if len(text) != 3 {
			err := errors.New(fewArguments)
			bot.SendError(err)
			return err
		}
		if err := bot.Backend.PutNewPage(text[1], text[2]); err != nil {
			bot.SendError(err)
			return err
		}
		mess.Text = dbUpdated
		bot.SendMessage(mess)
	default:
		switch ToLower(userText) {
		case "/start", reserv["start"]:
			mess.Text = OnStart()
		case "/помощь", "/help", reserv["help"]:
			mess.Text = onHelp
		case reserv["services"]:
			buff := serviceHelpMessage
			_, names, err := bot.Backend.GetPageUrlsAndNames()
			if err != nil {
				bot.SendError(err)
				return err
			}

			for i, n := range names {
				buff += fmt.Sprintf("\t%d. %s\n", i+1, strings.ToTitle(n))
			}
			mess.Text = buff
			if err := bot.Backend.PutCommand(GetHash(from), reserv["services"]); err != nil {
				bot.SendError(err)
				return err
			}
		case reserv["support"]:
			mess.Text = supportHelpMessage
			if err := bot.Backend.PutCommand(GetHash(from), reserv["support"]); err != nil {
				bot.SendError(err)
				return err
			}
		case reserv["search"]:
			mess.Text = searchHelpMessage
			if err := bot.Backend.PutCommand(GetHash(from), reserv["search"]); err != nil {
				bot.SendError(err)
				return err
			}
		case reserv["refresh"]:
			if err := bot.Backend.PutCommand(GetHash(from), reserv["refresh"]); err != nil {
				bot.SendError(err)
				return err
			}
			mess.Text = "Enter password"
		case reserv["addservice"]:
			if err := bot.Backend.PutCommand(GetHash(from), reserv["addservice"]); err != nil {
				bot.SendError(err)
				return err
			}
			mess.Text = "password|name|url"
		default:
			mess.Text = notFoundCommand
		}
		bot.SendMessage(mess)
	}

	return nil
}
