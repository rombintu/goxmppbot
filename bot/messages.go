package bot

import (
	"fmt"
	"regexp"
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
	return submess + ", напишите /help или помощь. Я вам помогу разобраться"
}

// Action on /help
func OnHelp() string {
	return `
Сервисы [сокращенно 'с'] - Вывести ссылки на ответы по сервисам
Поддержка - Помогу написать письмо в поддержку
Поиск - Помогу найти колег
	`
}

// Validate message for support mail
func ValidateSupport(message string) bool {
	data := strings.Split(message, ":")
	if len(data) < 2 {
		return false
	}
	matched, err := regexp.MatchString(`^[п|П]оддержка:[А-ЯЁ]+ [a-яА-ЯёЁ ]*`, message)
	if err != nil {
		return false
	}
	return matched
}

// Validate message for support mail
func ValidateSearch(message string) bool {
	data := strings.Split(message, ":")
	if len(data) < 2 {
		return false
	}
	matched, err := regexp.MatchString(`^[п|П]оиск:[\wа-яА-ЯЁ\s]+`, message)
	if err != nil {
		return false
	}
	return matched
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

func TrimS(message string) string {
	return strings.Trim(message, " ")
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
			forSupport := ValidateSupport(userText)
			forSearch := ValidateSearch(userText)
			if forSupport {
				dryData := strings.Split(userText, ":")
				emailData := ParseSubjectAndBody(dryData)
				resp, err := bot.OnSupport(from, emailData[0], emailData[1])
				if err != nil {
					bot.Logger.Error(err)
					mess.Text = "Произошла внутренняя ошибка: " + err.Error()
					continue
				}
				mess.Text = resp
				bot.SendMessage(mess)
				continue
			}

			if forSearch {
				mess.Text = "Выполняется..."
				bot.SendMessage(mess)
				dryData := strings.Split(userText, ":")
				count := "5"

				data := strings.ToLower(TrimS(dryData[1]))

				if len(dryData) > 2 {
					count = TrimS(dryData[2])
				}

				resp, err := GetUserByRegex(data, bot.Config.Contacts.Url, count)
				if err != nil {
					bot.Logger.Error(err)
					mess.Text = "Произошла внутренняя ошибка: " + err.Error()
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
			}

			switch strings.ToLower(userText) {
			case "/start", "start", "старт":
				mess.Text = OnStart()
			case "/помощь", "помощь", "п", "help", "/help":
				mess.Text = OnHelp()
			case "list", "лист", "сервисы", "сервис", "с":
				buff := ""
				for key, value := range bot.Config.Links {
					buff = buff + fmt.Sprintf("%s [%s]\n", value, key)
				}
				mess.Text = buff
			case "поддержка":
				mess.Text = `
Напишите свое обращение такого вида:
Поддержка:НАЗВАНИЕ_СЕРВИСА письмо
	Пример: Поддержка:СУДИС Все сломалось, помогите
			`
			case "поиск":
				mess.Text = `
Напишите свое обращение такого вида:
Поиск: ФИО_ПОЧТА_ДОЛЖНОСТЬ_КОМПАНИЯ
Примечание: Можно использовать регулярные выражения
Примечание: Добавьте в конце *: N*, чтобы регулировать выборку
	1. Пример: *Поиск: Иванов*
	2. Пример: *Поиск: Иванов: 10*
	3. Пример: *Поиск: ivanov*
			`
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
