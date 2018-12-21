package tokenapi

type VerifyTokenRequest struct {
	Token       string
	AppID       string
	TokenBody   string
	IsCheckBody bool //是否需要验证Body是否一致
}
