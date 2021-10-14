package bot

import (
	"fmt"
	"net/smtp"
)

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
		msg,
	); err != nil {
		return err
	}
	return nil
}
