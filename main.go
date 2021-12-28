package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	xmppbot "github.com/rombintu/goxmppbot/bot"
)

func main() {
	bot := xmppbot.NewBot()
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
	bot.ConfigurePlugins()
	bot.Plugins.Zabbix.GetToken()
	bot.Logger.Info("Bot token - OK")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitCh
		fmt.Println("Exit with 0")
		os.Exit(0)
	}()
	bot.Logger.Info("Bot started")
	bot.HandleMessage()
}
