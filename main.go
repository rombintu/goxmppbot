package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func main() {
	cfgPath := flag.String("config", "./config.toml", "Path to config.toml")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	bot := xmppbot.NewBot(*cfgPath)
	if err := bot.ConfigureBackand(); err != nil {
		bot.Logger.Error(err)
		os.Exit(1)
	}

	bot.Logger.Info("New bot was created")
	if err := bot.Connect(); err != nil {
		bot.Logger.Error(err)
		os.Exit(1)
	}
	bot.Logger.Info("Bot connected")

	// Plugins
	for _, plugin := range bot.Config.Default.Plugins {
		if plugin == "zabbix" {
			bot.ConfigurePlugins()
			bot.Plugins.Zabbix.GetToken()
			bot.Logger.Info("Bot token - OK")
		}
	}

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitCh
		fmt.Println("Exit with 0")
		bot.CloseLogFileOs()
		os.Exit(0)
	}()
	bot.Logger.Info("Bot started")
	bot.HandleMessage()
}
