package bot

import (
	"fmt"
	"os"

	zabbix "github.com/nixys/nxs-go-zabbix/v5"
)

func zabbixLogin(z *zabbix.Context, zbxHost, zbxUsername, zbxPassword string) {

	if err := z.Login(zbxHost, zbxUsername, zbxPassword); err != nil {
		fmt.Println("Login error:", err)
		os.Exit(1)
	} else {
		fmt.Println("Login: success")
	}
}

func zabbixLogout(z *zabbix.Context) {

	if err := z.Logout(); err != nil {
		fmt.Println("Logout error:", err)
		os.Exit(1)
	} else {
		fmt.Println("Logout: success")
	}
}

func Integration(zbxHost, zbxUsername, zbxPassword string) {
	var z zabbix.Context

	if zbxHost == "" || zbxUsername == "" || zbxPassword == "" {
		fmt.Println("Login error: host, username or password is empty")
		os.Exit(1)
	}

	/* Login to Zabbix server */
	zabbixLogin(&z, zbxHost, zbxUsername, zbxPassword)
	defer zabbixLogout(&z)

	/* Get all hosts */
	hObjects, _, err := z.HostGet(zabbix.HostGetParams{
		GetParameters: zabbix.GetParameters{
			Output: zabbix.SelectExtendedOutput,
		},
	})
	if err != nil {
		fmt.Println("Hosts get error:", err)
		return
	}

	/* Print names of retrieved hosts */
	fmt.Println("Hosts list:")
	for _, h := range hObjects {
		fmt.Println("-", h.Host)
	}
	// history, i, err := z.HistoryGet(
	// 	zabbix.HistoryGetParams{
	// 		History: 1,
	// 	},
	// )
	// fmt.Println("---------------------------------")
	// if err != nil {
	// 	fmt.Println("History get error:", err)
	// 	return
	// }
	// fmt.Println(history, i)
}
