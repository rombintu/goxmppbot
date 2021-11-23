package bot

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"
)

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
Сервисы - Вывести ссылки на ответы по сервисам
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
func ParseSubjectAndBody(message string) ([]string, error) {
	// support:subject body
	var subject, body string
	inner := strings.Split(message, " ")
	if len(inner) < 2 {
		return []string{}, errors.New("Несоответсвие шаблону")
	}
	subject = inner[0]
	body = strings.Join(inner[1:], " ")
	return []string{subject, body}, nil
}

func TrimS(message string) string {
	return strings.Trim(message, " ")
}

func ToLower(text string) string {
	return strings.ToLower(text)
}

func ToError(err error) string {
	return "Произошла внутренняя ошибка: " + err.Error()
}

func GetHash(login string) string {
	h := sha1.New()
	h.Write([]byte(login))
	return hex.EncodeToString(h.Sum(nil))
}
