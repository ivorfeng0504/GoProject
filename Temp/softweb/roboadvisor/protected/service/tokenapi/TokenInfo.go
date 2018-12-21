package tokenapi

type TokenInfo struct {
	Token       string
	AppID       string
	TokenBody   string
	LifeSeconds int //有效时间，单位为秒，如果输入不合法，默认为1800秒
}
