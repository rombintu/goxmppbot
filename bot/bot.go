package bot

import (
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	xmpp "github.com/mattn/go-xmpp"
	zabbixapi "github.com/rombintu/zabbix-api"
	"github.com/sirupsen/logrus"
)

type Default struct {
	Host          string        `toml:"Host"`
	Login         string        `toml:"Login"`
	Password      string        `toml:"Password"`
	DebugLevel    string        `toml:"DebugLevel"`
	DebugON       bool          `toml:"DebugOn"`
	LogFile       string        `toml:"LogFile"`
	UpdateChunk   time.Duration `toml:"UpdateChunk"`
	RefreshSecret string        `toml:"RefreshSecret"`
	Plugins       []string      `toml:"Plugins"`
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

// Struct for Zabbix-api
type ZabbixConf struct {
	Host string `toml:"Host"`
	User string `toml:"User"`
	Pass string `toml:"Password"`
}

// Struct BOT
type Bot struct {
	Client       *xmpp.Client
	Config       *Config
	Logger       *logrus.Logger
	LogFileOS    *os.File
	Backend      *Backend
	LastCommands map[string]string
	Plugins      Plugins
}

type Plugins struct {
	Zabbix zabbixapi.Zabbix
}

// Return configuration
func GetConfig(cfgPath string) *Config {
	confFile, err := os.ReadFile(cfgPath)
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
func NewBot(cfgPath string) *Bot {
	return &Bot{
		Config:       GetConfig(cfgPath),
		Logger:       logrus.New(),
		LastCommands: make(map[string]string),
	}
}

func (bot *Bot) OpenLogFileOs() error {
	f, err := os.OpenFile(bot.Config.Default.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	bot.LogFileOS = f
	return nil
}

func (bot *Bot) CloseLogFileOs() error {
	return bot.LogFileOS.Close()
}

func (bot *Bot) configureLogger() error {
	level, err := logrus.ParseLevel(bot.Config.Default.DebugLevel)
	if err != nil {
		return err
	}
	bot.Logger.SetLevel(level)
	bot.Logger.SetFormatter(&logrus.JSONFormatter{})
	bot.Logger.SetOutput(bot.LogFileOS)
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
	bot.Logger.Info("Bot plugins enabled")
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
	bot.Logger.Info("Bot client reconnected")
	return nil
}

// Try connect to server
func (bot *Bot) Connect() error {
	if err := bot.OpenLogFileOs(); err != nil {
		return err
	}
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
