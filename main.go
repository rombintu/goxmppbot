package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	xmppbot "github.com/rombintu/goxmppbot/bot"
)

type Default struct {
	Host     string // `toml:HOST`
	Login    string // `toml:LOGIN`
	Password string // `toml:PASSWORD`
	Debug    bool   // `toml:DEBUG`
}

type Config struct {
	Default Default
}

func GetConfig() Config {
	confFile, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("%v", err)
	}

	var conf Config

	if _, err := toml.Decode(string(confFile), &conf); err != nil {
		log.Fatalf("%v", err)
	}

	return conf
}

func main() {
	conf := GetConfig()
	bot := xmppbot.Bot{
		Host:     conf.Default.Host,
		Login:    conf.Default.Login,
		Password: conf.Default.Password,
		Debug:    conf.Default.Debug,
	}
	bot.Connect()

	err := bot.HandleMessage()
	if err != nil {
		log.Fatal(err)
	}
}
