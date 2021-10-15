package bot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
)

func encodeToBase64(msg []byte) []byte {
	var buff bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buff)
	encoder.Write(msg)
	defer encoder.Close()
	return buff.Bytes()
}

// func encodeToUTF8(msg string) []int {
// 	buff := make([]byte, len(msg))
// 	var encodeMess []int
// 	for _, r := range msg {
// 		encodeMess = append(encodeMess, utf8.EncodeRune(buff, r))
// 	}
// 	return encodeMess
// }

func (bot *Bot) SendToSupport(subject, body string) error {
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
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s.\r\n",
		bot.Config.Support.SupportEmail,
		subject,
		body,
	),
	)
	if err := smtp.SendMail(
		fmt.Sprintf("%s:%s", bot.Config.Support.Host, bot.Config.Support.Port),
		auth,
		fmt.Sprintf("%s@%s", bot.Config.Support.Login, bot.Config.Support.Host),
		to,
		encodeToBase64(msg),
	); err != nil {
		return err
	}
	return nil
}
