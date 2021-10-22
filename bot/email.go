package bot

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type Mail struct {
	MessageID string
	Subject   string
	From      string
	To        []string
	Body      []byte
}

func splitBytesToBlocks(body []byte) ([]byte, error) {
	return body, nil
}

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
	body, err := splitBytesToBlocks(wrapBase64([]byte(mail.Body)))
	if err != nil {
		log.Println(err)
	}
	buf.Write(body)

	return buf.Bytes()
}

func splitStringsToBlocks(body string) ([]byte, error) {
	n := 72 // 72 symbols
	b := new(bytes.Buffer)
	r := bufio.NewReader(b)
	buf := make([]byte, 0, n)
	for {
		n, err := io.ReadFull(r, buf[:cap(buf)])
		buf = buf[:n]
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != io.ErrUnexpectedEOF {
				fmt.Fprintln(os.Stderr, err)
				break
			}
		}
		prefix := []byte("\r\n")
		buf = append(buf, prefix...)
	}
	prefix := []byte("\r\n")
	buf = append(buf, prefix...)
	return buf, nil
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
	// splitBody, err := splitStringsToBlocks(body + fmt.Sprintf("\nСообщение отправлено от: %s", user))
	// if err != nil {
	// 	return err
	// }
	mail := Mail{
		MessageID: user,
		Subject:   subject,
		From:      fmt.Sprintf("%s@%s", bot.Config.Support.Login, bot.Config.Support.Host),
		To:        to,
		Body:      []byte(body + fmt.Sprintf("\n\nСообщение отправлено от: %s", user)),
		// Body:      splitBody,
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
