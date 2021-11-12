package bot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

type Mail struct {
	MessageID string
	Subject   string
	From      string
	To        []string
	Body      []byte
}

func splitBytesToBlocks(body []byte) []byte {
	var buf bytes.Buffer
	counter := 0
	for _, b := range body {
		buf.WriteByte(b)
		counter += 1
		if counter == 73 {
			buf.Write([]byte("\r\n"))
			counter = 0
		}
	}
	return buf.Bytes()
}

// func splitBytesToBlocks(body []byte) []byte {
// 	var buf bytes.Buffer
// 	countBlock := len(body)/72 + 1
// 	maxBlock := 72
// 	minBlock := 0
// 	for countBlock > 0 {
// 		if maxBlock > len(body)+countBlock*2 {
// 			buf.Write(body[minBlock:])
// 			break
// 		}
// 		buf.Write(body[:maxBlock])
// 		buf.Write([]byte("\r\n"))
// 		countBlock -= 1
// 		minBlock = maxBlock
// 		maxBlock += 75
// 	}
// 	return buf.Bytes()
// }

func wrapBase64(m []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(m)))
	base64.StdEncoding.Encode(b, m)
	return b
}

func BuildMail(mail Mail) []byte {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))

	subject := wrapBase64([]byte(mail.Subject))
	// buf.Write(subject)
	buf.WriteString(fmt.Sprintf("Subject: =?utf-8?b?%s?=\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString(fmt.Sprintf("Message-id: <%s>\r\n", mail.MessageID))

	// body := wrapBase64(mail.Body)
	body := splitBytesToBlocks(wrapBase64([]byte(mail.Body)))
	// if err != nil {
	// 	log.Println(err)
	// }
	buf.Write(body)

	return buf.Bytes()
}

func splitStringsToBlocks(body string) []byte {
	var buf bytes.Buffer
	counter := 0
	for _, r := range body {
		buf.WriteRune(r)
		counter += 1
		if counter == 73 {
			buf.Write([]byte("\r\n"))
			counter = 0
		}
	}
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
	// splitBody := splitStringsToBlocks(body + fmt.Sprintf("\nСообщение отправлено от: %s", user))
	// if err != nil {
	// 	return err
	// }
	mail := Mail{
		MessageID: user,
		Subject:   subject,
		From:      fmt.Sprintf("%s@%s", bot.Config.Support.Login, bot.Config.Support.Host),
		To:        to,
		Body:      []byte(body + fmt.Sprintf("\n\nСообщение отправлено от: %s", user)),
		// Body: splitBody,
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
