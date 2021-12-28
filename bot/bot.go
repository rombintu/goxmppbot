package bot

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	xmpp "github.com/mattn/go-xmpp"
	zabbixapi "github.com/rombintu/goxmppbot/plugins/zabbix-api"
	"github.com/sirupsen/logrus"
)

type Default struct {
	Host          string `toml:"Host"`
	Login         string `toml:"Login"`
	Password      string `toml:"Password"`
	DebugLevel    string `toml:"DebugLevel"`
	DebugON       bool   `toml:"DebugOn"`
	RefreshSecret string `toml:"RefreshSecret"`
}

type Support struct {
	Host         string `toml:"Host"`
	Port         string `toml:"Port"`
	Login        string `toml:"LoginWithoutHost"`
	Password     string `toml:"Password"`
	SupportEmail string `toml:"SupportEmail"`
}

type Contacts struct {
	Url string `toml:"URL"`
}

type DBConf struct {
	Master string `toml:"Master"`
	Slave  string `toml:"Slave"`
	Dev    bool   `toml:"Dev"`
	Multi  bool   `toml:"Multi"`
}

type Config struct {
	Default  Default    `toml:"DEFAULT"`
	Support  Support    `toml:"SUPPORT"`
	Contacts Contacts   `toml:"CONTACTS"`
	DBConf   DBConf     `toml:"DBCONF"`
	Zabbix   ZabbixConf `toml:"ZABBIX"`
}

// Struct BOT
type Bot struct {
	Client       *xmpp.Client
	Config       *Config
	Logger       *logrus.Logger
	Backend      *Backend
	LastCommands map[string]string
	Plugins      Plugins
}

type Plugins struct {
	Zabbix zabbixapi.Zabbix
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
		Config:       GetConfig(),
		Logger:       logrus.New(),
		LastCommands: make(map[string]string),
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
	bot.Backend = NewBackend(bot.Config.DBConf)
	if err := bot.Backend.Init(); err != nil {
		return err
	}
	return nil
}

func (bot *Bot) ConfigurePlugins() {
	// enable zabbix-api
	bot.Plugins.Zabbix = zabbixapi.NewZabbix(
		bot.Config.Zabbix.Host,
		bot.Config.Zabbix.User,
		bot.Config.Zabbix.Pass,
	)
	bot.Logger.Info("bot plugins enabled")
}

func (bot *Bot) ReConnnect() error {
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
	bot.Logger.Info("bot client reconnected")
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
