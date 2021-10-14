package bot

import "testing"

func TestSendToSupport(t *testing.T) {
	bot := NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	if err := bot.SendToSupport("themeTest", "bodyTest"); err != nil {
		t.Error(err)
	}
}
