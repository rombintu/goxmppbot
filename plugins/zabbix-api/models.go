package zabbixapi

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type DefaultJsonFields struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      int    `json:"id"`
}

type ProblemGetParams struct {
	Output     string   `json:"output"`
	Resent     string   `json:"resent"`
	SorteField []string `json:"sortfield"`
	Sortorder  string   `json:"sortorder"`
}

type ProblemGet struct {
	DefaultJsonFields
	Params ProblemGetParams `json:"params"`
	Auth   string           `json:"auth"`
}

type ProblemGetResp struct {
	Result []Problems `json:"result"`
	Error  Error      `json:"error"`
}

type Problems struct {
	Clock        string `json:"clock"`
	Text         string `json:"name"`
	Acknowledged string `json:"acknowledged"`
}

type UserLoginParams struct {
	Username string `json:"user"`
	Password string `json:"password"`
	UserData bool   `json:"userData"`
}

type UserLogin struct {
	DefaultJsonFields
	Params UserLoginParams `json:"params"`
}

type UserLoginResp struct {
	Result string `json:"result"`
	Error  Error  `json:"error"`
}

type ResponseWithAuth struct {
	DefaultJsonFields
	Params interface{} `json:"params"`
	Auth   string      `json:"auth"`
}

type ResponseWithoutAuth struct {
	DefaultJsonFields
	Params interface{} `json:"params"`
}
