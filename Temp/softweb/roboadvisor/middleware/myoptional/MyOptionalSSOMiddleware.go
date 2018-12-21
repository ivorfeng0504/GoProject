package middleware

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/contract/myoptional"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service/sso"
	"git.emoney.cn/softweb/roboadvisor/util/crypto"
	"git.emoney.cn/softweb/roboadvisor/web_front_myoptional/handlers/error"
	"github.com/devfeel/dotweb"
	"strings"
)

type MyOptionalSSOMiddleware struct {
	dotweb.BaseMiddlware
}

const (
	SSOHashCookieName = "ssohash"
)

func getSSOUrl(ctx dotweb.Context) (sso string, ssoHash string) {
	defer func() {
		if err := recover(); err != nil {
			sso = ""
			ssoHash = ""
		}
	}()

	ssoTemplate := "Version=%s&uid=%s&cid=%s&pid=%s&sid=%s&tid=%s&agentid=%s&clienttype=%s&OutOfDate=%s&pluglet=%s&token=%s&bata=%s"
	//rand := ctx.QueryString("rand")
	Version := ctx.QueryString("Version")
	uid := ctx.QueryString("uid")
	cid := ctx.QueryString("cid")
	pid := ctx.QueryString("pid")
	sid := ctx.QueryString("sid")
	tid := ctx.QueryString("tid")
	agentid := ctx.QueryString("agentid")
	clienttype := ctx.QueryString("clienttype")
	OutOfDate := ctx.QueryString("OutOfDate")
	pluglet := ctx.QueryString("pluglet")
	token := ctx.QueryString("token")
	bata := ctx.QueryString("bata")
	url := fmt.Sprintf(ssoTemplate, Version, uid, cid, pid, sid, tid, agentid, clienttype, OutOfDate, pluglet, token, bata)
	ssoHash = _crypto.MD5(url)[:8]
	return url, ssoHash
}

func getSSOHash(ctx dotweb.Context) (ssoHash string, err error) {
	ssoHash, err = ctx.ReadCookieValue(SSOHashCookieName)
	return ssoHash, err
}

func setSSOHash(ctx dotweb.Context, ssoHash string) {
	ctx.SetCookieValue(SSOHashCookieName, ssoHash, 0)
}

func (middleware *MyOptionalSSOMiddleware) Handle(ctx dotweb.Context) error {
	ssoUrl, ssoHash := getSSOUrl(ctx)
	//没有SSO直接返回未授权
	if len(ssoUrl) == 0 || len(ssoHash) == 0 {
		return _error.NotAuth(ctx)
	}
	user := contract_myoptional.UserInfo(ctx)
	oldSsoHash, err := getSSOHash(ctx)
	//如果用户已登录且SSO未发生变更则不进行SSO解析
	if err == nil && len(oldSsoHash) > 0 && user != nil && ssoHash == oldSsoHash {
		middleware.Next(ctx)
		return nil
	}
	ssoUser, err := sso.DecryptSso(ctx.Request().RawQuery(), ctx.RemoteIP())
	if err != nil {
		global.InnerLogger.Error(err, "SSO 解析异常")
		return _error.ServerError(ctx)
	}
	if ssoUser == nil {
		return _error.NotAuth(ctx)
	}
	account := ssoUser.Account
	uid := ssoUser.Uid
	//如果SSO中存在CID 则使用CID覆盖UID
	if ssoUser.Cid > 0 {
		uid = ssoUser.Cid
	}
	if len(account) == 0 && uid <= 0 {
		global.InnerLogger.DebugFormat("SSO 解析异常  Account与UID同时不存在 ssoUrl=%s", ssoUrl)
		return _error.NotAuth(ctx)
	}
	user = &live.User{}
	user.Account = account
	user.UID = uid

	nickName := account
	accountLen := len(account)
	if accountLen > 0 {
		if accountLen >= 3 {
			mask := "****"
			if accountLen < 7 {
				mask = strings.Repeat("*", accountLen-3)
			}
			nickName = account[:3] + mask
			if accountLen > 7 {
				start := accountLen - 3
				if start < 7 {
					start = 7
				}
				nickName += account[start:]
			}
		} else {
			nickName = account
		}
	}

	user.NickName = nickName

	//用户登录
	contract_myoptional.SetUserInfo(ctx, user)
	setSSOHash(ctx, ssoHash)
	middleware.Next(ctx)
	return nil
}

func NewMyOptionalSSOMiddleware() *MyOptionalSSOMiddleware {
	return &MyOptionalSSOMiddleware{}
}
