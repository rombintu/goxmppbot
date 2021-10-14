package bot

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	xmpp "github.com/mattn/go-xmpp"
	"github.com/sirupsen/logrus"
)

type Default struct {
	Host       string // `toml:HOST`
	Login      string // `toml:LOGIN`
	Password   string // `toml:PASSWORD`
	DebugLevel string
	DebugON    bool
}

type Support struct {
	Host         string
	Port         string
	Login        string
	Password     string
	SupportEmail string //`toml:SUPPORT_EMAIL`
}

type Config struct {
	Default Default
	Support Support
	Links   map[string]string
}

// Struct BOT
type Bot struct {
	Client *xmpp.Client
	Config *Config
	Logger *logrus.Logger
}

// Return configuration
func GetConfig() *Config {
	confFile, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("%v", err)
	}

	var conf Config

	if _, err := toml.Decode(string(confFile), &conf); err != nil {
		log.Fatalf("%v", err)
	}

	return &conf
}

// Return new bot
func NewBot() *Bot {
	return &Bot{
		Config: GetConfig(),
		Logger: logrus.New(),
	}
}

func (bot *Bot) configureLogger() error {
	level, err := logrus.ParseLevel(bot.Config.Default.DebugLevel)
	if err != nil {
		return err
	}

	bot.Logger.SetLevel(level)

	return nil
}

// Try connect to server
func (bot *Bot) Connect() error {
	err := bot.configureLogger()
	if err != nil {
		bot.Logger.Error("Error configure logger: ", err.Error())
		return err
	}
	client, err := xmpp.NewClientNoTLS(
		bot.Config.Default.Host,
		bot.Config.Default.Login,
		bot.Config.Default.Password,
		bot.Config.Default.DebugON,
	)
	if err != nil {
		bot.Logger.Error("Error connect: ", err.Error())
		return err
	}
	bot.Client = client
	return nil
}
