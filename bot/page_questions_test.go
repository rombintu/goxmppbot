// "http://it.mvd.ru/faq/38"

package bot_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rombintu/goxmppbot/bot"
)

func TestGetPage(t *testing.T) {
	data, err := bot.GetPage("http://it.mvd.ru/faq/38")
	if err != nil {
		t.Fatal(err)
	}
	var page bot.Page
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatal(err)
	}
	fmt.Println(page)
}
