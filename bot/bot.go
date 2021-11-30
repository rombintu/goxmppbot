package bot

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	xmpp "github.com/mattn/go-xmpp"
	"github.com/sirupsen/logrus"
)

type Default struct {
	Host          string `toml:"HOST"`
	Login         string `toml:"LOGIN"`
	Password      string `toml:"PASSWORD"`
	DebugLevel    string `toml:"DEBUGLEVEL"`
	DebugON       bool   `toml:"DEBUG"`
	RefreshSecret string `toml:"REFRESH_SECRET"`
}

type Support struct {
	Host         string `toml:"HOST"`
	Port         string `toml:"PORT"`
	Login        string `toml:"LOGIN"`
	Password     string `toml:"PASSWORD"`
	SupportEmail string `toml:"SUPPORTEMAIL"`
}

type Contacts struct {
	Url string `toml:"URL"`
}

type BackendConf struct {
	Host         string `toml:"HOST"`
	Port         string `toml:"PORT"`
	User         string `toml:"USER"`
	Password     string `toml:"PASSWORD"`
	SSLMode      string `toml:"SSLMODE"`
	DatabaseName string `toml:"DATABASE"`
	Dev          bool   `toml:"DEV"`
	Multi        bool   `toml:"MULTI"`
}

type Config struct {
	Default          Default           `toml:"DEFAULT"`
	Support          Support           `toml:"SUPPORT"`
	Contacts         Contacts          `toml:"CONTACTS"`
	BackendConf      BackendConf       `toml:"BACKENDCONF"`
	BackendConfSlave BackendConf       `toml:"BACKENDCONF_SLAVE"`
	Links            map[string]string `toml:"LINKS"`
}

// Struct BOT
type Bot struct {
	Client  *xmpp.Client
	Config  *Config
	Logger  *logrus.Logger
	Backend *Backend
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

func (bot *Bot) ConfigureBackand() error {
	bot.Backend = NewBackend(bot.Config.BackendConf, bot.Config.BackendConfSlave)
	if err := bot.Backend.Init(); err != nil {
		return err
	}
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
