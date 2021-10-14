package bot_test

import (
	"fmt"
	"testing"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func TestOnStart(t *testing.T) {
	fmt.Println(xmppbot.OnStart())
}
