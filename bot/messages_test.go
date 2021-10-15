package bot_test

import (
	"fmt"
	"testing"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func TestOnStart(t *testing.T) {
	fmt.Println(xmppbot.OnStart())
}

func TestValidateSupport(t *testing.T) {
	var tests = []struct {
		inter string
		want  bool
	}{
		{"поддержка:СЭП помогите", true},
		{"Поддержка:СУДИС помогите мне плис", true},
		{"поддержка:СУДИС хелп хелп хелп", true},
		{"Поддержка:СЭП у меня лапки", true},

		{"поддержка: СЭП помогите", false},
		{"подержка:СЭП помогите", false},
		{"поддержка:СЭП", false},
		{"поддержка:сэп помогите мне", false},
		{" поддержка:СЭП помогите", false},
	}

	for _, test := range tests {
		res := fmt.Sprintf("case (%s)", test.inter)
		t.Run(res, func(t *testing.T) {
			got := xmppbot.ValidateSupport(test.inter)
			if got != test.want {
				t.Errorf("GOT: [%t] WANT: [%t]", got, test.want)
			}
		})
	}
}
