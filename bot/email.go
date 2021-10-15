package bot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

type Mail struct {
	Subject string
	From    string
	To      []string
	Body    []byte
}

func BuildMail(mail Mail) []byte {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))

	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")

	b := make([]byte, base64.StdEncoding.EncodedLen(len(mail.Body)))
	base64.StdEncoding.Encode(b, mail.Body)
	buf.Write(b)

	buf.WriteString("--")

	return buf.Bytes()
}

func (bot *Bot) SendToSupport(user, subject, body string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		bot.Config.Support.Login,
		bot.Config.Support.Password,
		bot.Config.Support.Host,
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{bot.Config.Support.SupportEmail}

	mail := Mail{
		Subject: subject,
		From:    fmt.Sprintf("%s@%s", bot.Config.Support.Login, bot.Config.Support.Host),
		To:      to,
		Body:    []byte(body + fmt.Sprintf("\n\nСообщение отправлено от: %s", user)),
	}

	message := BuildMail(mail)

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%s", bot.Config.Support.Host, bot.Config.Support.Port),
		auth,
		fmt.Sprintf("%s@%s", bot.Config.Support.Login, bot.Config.Support.Host),
		to,
		message,
	); err != nil {
		return err
	}
	return nil
}
