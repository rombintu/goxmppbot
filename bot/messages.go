package bot

import (
	"encoding/json"
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

// Action on /support. Send mail to support
func (bot *Bot) OnSupport(user, subject, body string) (string, error) {
	if err := bot.SendToSupport(user, subject, body); err != nil {
		bot.Logger.Error(err)
		return "", err
	}
	return "Ваша заявка отправлена в обработку", nil
}

func (bot *Bot) SendError(mess xmpp.Chat, err error) {
	mess.Text = err.Error()
	bot.SendMessage(mess)
}

// Loop func, listening command from users
func (bot *Bot) HandleMessage() {
	mainCh := make(chan interface{})

	go bot.Handler(mainCh)
	for {
		data, err := bot.Client.Recv()
		if err != nil {
			bot.Logger.Error(err)
		}
		mainCh <- data
	}

}

func (bot *Bot) Handler(c chan interface{}) {
	timeCh := time.NewTicker(bot.Config.Default.UpdateChunk * time.Minute)
	for {
		select {
		case data, open := <-c:
			if !open {
				break
			}
			switch data.(type) {
			case xmpp.Chat:
				if err := bot.Run(data); err != nil {
					bot.Logger.Error(err)
					if err := bot.ReConnnect(); err != nil {
						bot.Logger.Error(err)
					}
				}
			}
		case <-timeCh.C:
			bot.LastCommands = make(map[string]string)
			if err := bot.ReConnnect(); err != nil {
				bot.Logger.Error(err)
			}
		}
	}
}

func (bot *Bot) Run(data interface{}) error {
	if data.(xmpp.Chat).Text == "" || data.(xmpp.Chat).Text == " " {
		return nil
	}
	bot.Logger.Info(data)
	mess := CreateMessage()
	from := data.(xmpp.Chat).Remote
	mess.Remote = from
	// mess.Subject = "bothelper"
	userText := data.(xmpp.Chat).Text

	lc := bot.LastCommands[GetHash(from)]
	bot.LastCommands[GetHash(from)] = ""

	switch lc {
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
			bot.SendError(mess, err)
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
			bot.SendError(mess, err)
			return err
		}
		resp, err := bot.OnSupport(from, emailData[0], emailData[1])
		if err != nil {
			bot.SendError(mess, err)
			return err
		}
		mess.Text = resp
		bot.SendMessage(mess)
	case reserv["services"]:
		data, err := bot.Backend.GetJsonByName(ToLower(userText))
		if err != nil {
			bot.SendError(mess, err)
			return err
		}

		if len(data) == 0 {
			mess.Text = notFound
			bot.SendMessage(mess)
			return nil
		}
		var page Page
		if err := json.Unmarshal(data, &page); err != nil {
			bot.SendError(mess, err)
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
			bot.SendError(mess, err)
			return err
		}
		for _, u := range urls {
			page, err := GetPage(u)
			if err != nil {
				bot.SendError(mess, err)
				return err
			}
			if err := bot.Backend.UpdatePage(page, u); err != nil {
				bot.SendError(mess, err)
				return err
			}
		}
		mess.Text = dbUpdated
		bot.SendMessage(mess)
	case reserv["addservice"]:
		text := strings.Split(userText, "|")
		if text[0] != bot.Config.Default.RefreshSecret || len(text) != 3 {
			return nil
		}
		mess.Text = loading
		bot.SendMessage(mess)

		if err := bot.Backend.PutNewPage(text[1], text[2]); err != nil {
			bot.SendError(mess, err)
			return err
		}
		mess.Text = dbUpdated
		bot.SendMessage(mess)
	default:
		switch ToLower(userText) {
		case "/start", reserv["start"]:
			mess.Text = OnStart()
		case "/помощь", "help", reserv["help"]:
			mess.Text = onHelp
		case reserv["services"]:
			buff := serviceHelpMessage
			_, names, err := bot.Backend.GetPageUrlsAndNames()
			if err != nil {
				bot.SendError(mess, err)
				return err
			}

			if len(names) == 0 {
				buff += "\t[Сервисы не найдены в базе данных]"
			} else {
				for i, n := range names {
					buff += fmt.Sprintf("\t%d. %s\n", i+1, strings.ToTitle(n))
				}
			}

			mess.Text = buff

		case reserv["support"]:
			mess.Text = supportHelpMessage
		case reserv["search"]:
			mess.Text = searchHelpMessage
		case reserv["refresh"]:
			mess.Text = "Enter password"
		case reserv["addservice"]:
			mess.Text = "password|name|url"
		case reserv["zabbix"], "заббикс", "забикс":
			buff := ""
			on := false
			for _, plugin := range bot.Config.Default.Plugins {
				if plugin == "zabbix" {
					on = true
				}

			}
			if !on {
				bot.SendError(mess, errors.New(pluginNotEnabled))
				return nil
			}

			problems, err := bot.Plugins.Zabbix.GetProblems()
			if err != nil || problems.Error.Message != "" {
				bot.SendError(mess, err)
				bot.SendError(mess, errors.New(problems.Error.Message))
				return err
			}
			for i, p := range problems.Result {
				ack := "Нет"
				if p.Acknowledged == "1" {
					ack = "Да"
				}

				clock, _ := StrToTime(p.Clock)

				buff += "--------\n"
				buff += fmt.Sprintf(
					"%d. [\n\tВремя: %s\n\tСообщение: %s\n\tРешено: %s\n]\n",
					i+1, clock, p.Text, ack)
			}
			mess.Text = buff
		default:
			mess.Text = notFoundCommand
		}
		bot.LastCommands[GetHash(from)] = ToLower(userText)
		bot.SendMessage(mess)
	}

	return nil
}
