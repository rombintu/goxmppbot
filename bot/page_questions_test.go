// "http://it.mvd.ru/faq/38"

package bot_test

import (
	"fmt"
	"testing"

	"github.com/rombintu/goxmppbot/bot"
)

func TestGetPage(t *testing.T) {
	data, err := bot.GetPage("http://it.mvd.ru/faq/39")
	if err != nil {
		t.Fatal(err)
	}
	// var page bot.Page
	// if err := json.Unmarshal(data, &page); err != nil {
	// 	t.Fatal(err)
	// }
	fmt.Println(data)
}

// func TestGetData(t *testing.T) {
// 	bot := bot.NewBot()
// 	if err := bot.ConfigureBackand(); err != nil {
// 		bot.Logger.Error(err)
// 		t.Fatal(err)
// 	}
// 	_, _, data, err := bot.Backend.GetPageUrlsAndNames()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(data)
// }
