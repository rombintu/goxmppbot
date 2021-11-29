package bot

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Question struct {
	Subquestion []string
	Subanswer   []string
}

type Page struct {
	Questions []Question
}

func GetPage(url string) ([]byte, error) {
	// Perform request
	resp, err := http.Get(url)
	if err != nil {
		print(err)
		return []byte{}, err
	}
	// Cleanup when this function ends
	defer resp.Body.Close()
	// Read & parse response data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	var quests string
	// var quests []string
	// Print content of <title></title>
	doc.Find(".faq-list").Each(func(i int, s *goquery.Selection) {
		quests = strings.TrimSpace(s.Text())
	})

	arr := strings.Split(quests, "развернуть")
	var questions []Question
	var quest, answ []string
	for _, a := range arr {
		tmp := strings.Split(strings.TrimSpace(a), "\n")
		quest = append(quest, tmp[0])
		answ = append(answ, strings.Join(tmp[1:], " "))
		questions = append(questions, Question{
			Subquestion: quest,
			Subanswer:   answ,
		})
	}
	page := Page{
		Questions: questions,
	}
	data, err := json.Marshal(page)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
