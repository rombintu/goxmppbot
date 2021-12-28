package zabbixapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ZabbixConf struct {
	Host string `toml:"Host"`
	User string `toml:"User"`
	Pass string `toml:"Password"`
}

type Zabbix struct {
	Conf  ZabbixConf
	Token string
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

func (z *Zabbix) GetRequest(respBody *bytes.Buffer) ([]byte, error) {

	resp, err := http.Post(fmt.Sprintf("http://%s/zabbix/api_jsonrpc.php", z.Conf.Host), "application/json-rpc", respBody)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (z *Zabbix) WorkJson(newStruct interface{}) ([]byte, error) {
	postBody, err := json.Marshal(newStruct)
	if err != nil {
		return []byte{}, err
	}
	responseBody := bytes.NewBuffer(postBody)
	data, err := z.GetRequest(responseBody)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (z *Zabbix) BuildJson(params interface{}, method string, id int) ([]byte, error) {
	newStruct := ResponseWithoutAuth{
		DefaultJsonFields: DefaultJsonFields{
			Jsonrpc: "2.0",
			Method:  method,
			Id:      id,
		},
		Params: params,
	}
	data, err := z.WorkJson(newStruct)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (z *Zabbix) BuildJsonWithAuth(params interface{}, method string, id int) ([]byte, error) {
	newStruct := ResponseWithAuth{
		DefaultJsonFields: DefaultJsonFields{
			Jsonrpc: "2.0",
			Method:  method,
			Id:      id,
		},
		Params: params,
		Auth:   z.Token,
	}
	data, err := z.WorkJson(newStruct)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (z *Zabbix) GetToken() (UserLoginResp, error) {

	params := UserLoginParams{
		Username: z.Conf.User,
		Password: z.Conf.Pass,
		UserData: false,
	}

	data, err := z.BuildJson(params, "user.login", 1)
	if err != nil {
		return UserLoginResp{}, err
	}
	var result UserLoginResp
	if err := json.Unmarshal(data, &result); err != nil {
		return UserLoginResp{}, err
	}
	if result.Result != "" {
		z.Token = result.Result
	} else {
		return result, errors.New(result.Error.Message)
	}
	return result, nil
}

func StrToTime(timeStr string) (string, error) {
	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return "Error", err
	}
	tm := time.Unix(i, 0).String()
	return tm, nil
}

func (z *Zabbix) GetProblems() (ProblemGetResp, error) {
	params := ProblemGetParams{
		Output:     "extend",
		Resent:     "true",
		SorteField: []string{"eventid"},
		Sortorder:  "DESC",
	}

	data, err := z.BuildJsonWithAuth(params, "problem.get", 1)
	if err != nil {
		return ProblemGetResp{}, err
	}
	var result ProblemGetResp
	if err := json.Unmarshal(data, &result); err != nil {
		return ProblemGetResp{}, err
	}
	return result, nil
}
