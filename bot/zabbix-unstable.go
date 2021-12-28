package bot

import (
	"errors"
	"fmt"

	zabbix "github.com/nixys/nxs-go-zabbix/v5"
)

type ZabbixConf struct {
	Host string `toml:"Host"`
	User string `toml:"User"`
	Pass string `toml:"Password"`
}

type Zabbix struct {
	Conf    ZabbixConf
	Conn    zabbix.Context
	Hosts   []zabbix.HostObject
	Actions []zabbix.ActionObject
	// History *[]zabbix.HistoryLogObject
	History interface{}
}

func NewZabbix(host, user, pass string) Zabbix {
	c := ZabbixConf{
		Host: host,
		User: user,
		Pass: pass,
	}
	return Zabbix{
		Conf: c,
	}
}

// Get all hosts
func (z *Zabbix) getHosts() error {
	hObjects, _, err := z.Conn.HostGet(zabbix.HostGetParams{
		GetParameters: zabbix.GetParameters{
			Output: zabbix.SelectExtendedOutput,
		},
	})
	if err != nil {
		return err
	}
	z.Hosts = hObjects
	return nil
}

// Обновляет список хостов
func (z *Zabbix) GetHosts() error {
	if err := z.Connect(); err != nil {
		return err
	}
	defer z.Conn.Logout()

	if err := z.getHosts(); err != nil {
		return err
	}
	return nil
}

func (z *Zabbix) GetActions() error {
	if err := z.Connect(); err != nil {
		return err
	}
	defer z.Conn.Logout()

	actions, _, err := z.Conn.ActionGet(
		zabbix.ActionGetParams{
			SelectFilter:             "extend",
			SelectOperations:         "extend",
			SelectRecoveryOperations: "extend",
		},
	)
	if err != nil {
		return err
	}

	z.Actions = actions

	return nil
}

func (z *Zabbix) GetHistory() error {
	if err := z.Connect(); err != nil {
		return err
	}
	defer z.Conn.Logout()
	history, _, err := z.Conn.HistoryGet(
		zabbix.HistoryGetParams{
			History: 1,
		},
	)
	if err != nil {
		return err
	}

	// z.History = history.(*[]zabbix.HistoryTextObject)
	// z.History = history.(*[]zabbix.HistoryLogObject)
	z.History = history

	return nil
}

func (z *Zabbix) Connect() error {

	if z.Conf.Host == "" || z.Conf.User == "" || z.Conf.Pass == "" {
		return errors.New(fmt.Sprintln("Login error: host, username or password is empty"))
	}

	/* Login to Zabbix server */
	if err := z.Conn.Login(z.Conf.Host, z.Conf.User, z.Conf.Pass); err != nil {
		return err
	}

	return nil
}
