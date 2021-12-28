package main_test

import (
	"fmt"
	"testing"
	"time"

	xmppbot "github.com/rombintu/goxmppbot/bot"
	zabbixapi "github.com/rombintu/zabbix-api"
)

func TestParseConfig(t *testing.T) {
	bot := xmppbot.NewBot()

	fmt.Println(bot.Config)
}

func TestSendToSupport(t *testing.T) {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	if err := bot.SendToSupport("test@test.ru", "СУДИС", "тело текста -- тест test $#$FDF"); err != nil {
		t.Error(err)
	}
}
func TestSendToSupportManySymbols(t *testing.T) {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	if err := bot.SendToSupport("test@test.ru", "СУДИС", "gsdfffffffffffffffffffffffgdfsgdfgsdfgdsfgsdfgfdshfghfdggggggggggggggggggggggggggggggggggggggggggggggghgsfhsfghfghdfhfghfg"); err != nil {
		t.Error(err)
	}
}

func TestSendOOB(t *testing.T) {
	bot := xmppbot.NewBot()
	err := bot.Connect()
	if err != nil {
		t.Error(err)
	}
	mess := xmppbot.CreateMessage()
	mess.Text = "test"
	mess.Remote = ""
	mess.Subject = "bothelper"
	mess.Ooburl = "http://risovach.ru/upload/2014/08/mem/nu-pochemu_59642579_orig_.jpg"
	mess.Stamp.Add(10 * time.Second)
	bot.SendMessage(mess)
}

func TestGetXslxFromUrl(t *testing.T) {
	data, err := xmppbot.GetXslxFromUrl("http://0.0.0.0:8000/ЗНИ_СЭП_приклад.xlsx")
	if err != nil {
		t.Error(err)
	}
	rmap, err := xmppbot.OpenXslxFile(data)
	if err != nil {
		t.Error(err)
	}
	for k, v := range rmap {
		fmt.Println(k, ":", v)
	}
}

// func TestZabbixConn(t *testing.T) {
// 	z := xmppbot.NewZabbix("http://192.168.213.127/zabbix/api_jsonrpc.php", "*", "*")
// 	z.GetHosts()
// 	// z.GetHistory()
// 	// z.GetActions()
// 	// for _, h := range z.Hosts {
// 	// 	fmt.Println(h.Host, h.HostID, h.)
// 	// }
// 	fmt.Println("-----------")
// 	// for _, hist := range *z.History {
// 	// 	fmt.Println(hist.Value)
// 	// }
// 	// fmt.Println(z.History)
// 	fmt.Printf("%+v", z.Hosts)
// }

func TestGetTokenZabbix(t *testing.T) {
	z := zabbixapi.NewZabbix("192.168.213.127", "Admin", "zabbix")
	_, err := z.GetToken()
	if err != nil {
		t.Fatal(err)
	}
	problems, err := z.GetProblems()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(problems.Result)
	fmt.Println(problems.Error)
}
