package middleware

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"emoney.cn/fundchannel/util/crypto"
	"emoney.cn/fundchannel/contract"
	"emoney.cn/fundchannel/web_fundchannel/handlers/error"
	"emoney.cn/fundchannel/protected/service/sso"
	"emoney.cn/fundchannel/global"
	"github.com/devfeel/mapper"
	_url "net/url"
	"emoney.cn/fundchannel/protected/model"
)

type SSOMiddleware struct {
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

	ssoTemplate := "Version=%s&uid=%s&pid=%s&sid=%s&tid=%s&agentid=%s&clienttype=%s&OutOfDate=%s&pluglet=%s&cid=%s&token=%s&bata=%s"
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
	url := fmt.Sprintf(ssoTemplate, Version, uid, pid, sid, tid, agentid, clienttype, OutOfDate, pluglet, cid, token, bata)

	ssoQuery := fmt.Sprintf(ssoTemplate, Version, uid, pid, sid, tid, agentid, clienttype, OutOfDate, pluglet, cid, _url.QueryEscape(token), bata)
	contract.SetSSOQuery(ctx, ssoQuery)
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

func (middleware *SSOMiddleware) Handle(ctx dotweb.Context) error {
	ssoUrl, ssoHash := getSSOUrl(ctx)
	//没有SSO直接返回未授权
	if len(ssoUrl) == 0 || len(ssoHash) == 0 {
		return _error.NotAuth(ctx)
	}

	user := contract.UserInfo(ctx)

	oldSsoHash, err := getSSOHash(ctx)

	//如果用户已登录且SSO未发生变更则不进行SSO解析
	if err == nil && len(oldSsoHash) > 0 && user != nil && ssoHash == oldSsoHash {
		middleware.Next(ctx)
		return nil
	}
	//用户身份发生切换或者用户未登录 重置用户身份
	user = new(fund.User)

	ssoUser, err := sso.DecryptSso(ctx.Request().RawQuery(), ctx.RemoteIP())
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "SSO 解析异常 ssoUrl=%s", ssoUrl)
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
	err = mapper.Mapper(ssoUser, user)
	if err != nil {
		global.InnerLogger.ErrorFormat(err, "SSO 解析 用户信息转换异常 ssoUrl=%s", ssoUrl)
		return _error.ServerError(ctx)
	}
	//用户登录
	contract.SetUserInfo(ctx, user)
	setSSOHash(ctx, ssoHash)
	middleware.Next(ctx)
	return nil
}

func NewSSOMiddleware() *SSOMiddleware {
	return &SSOMiddleware{}
}
