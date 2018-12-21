package middleware

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/util/crypto"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/error"
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract/expert"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/sso"
	"git.emoney.cn/softweb/roboadvisor/protected/model/expert"
	"net/url"
)

type ExpertSSOMiddleware struct {
	dotweb.BaseMiddlware
}

func getSSOUrl(ctx dotweb.Context) (sso string, ssoHash string) {
	defer func() {
		if err := recover(); err != nil || ctx.Request().RawQuery()=="" {
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
	url := fmt.Sprintf(ssoTemplate, Version, uid, cid, pid, sid, tid, agentid, clienttype, OutOfDate, pluglet, url.QueryEscape(token), bata)
	ssoHash = _crypto.MD5(url)[:8]
	return url, ssoHash
}

func getSSOHash(ctx dotweb.Context) (ssoHash string, err error) {
	ssoHash, err = ctx.ReadCookieValue("expertssohash")
	return ssoHash, err
}

func setSSOHash(ctx dotweb.Context, ssoHash string) {
	ctx.SetCookieValue("expertssohash", ssoHash, 0)
}


func setSSOUrl(ctx dotweb.Context, cid int64,ssoUrl string) {
	ctx.SetCookieValue("expertnews.uid",strconv.FormatInt(cid,10) ,0)
	ctx.SetCookieValue("expertnews.focusssourl",ssoUrl,0)
}

func (middleware *ExpertSSOMiddleware) Handle(ctx dotweb.Context) error {
	ssoUrl, ssoHash := getSSOUrl(ctx)
	loginuser := contract_expert.ExpertUserInfo(ctx)

	//没有SSO直接返回未授权
	if loginuser == nil && (len(ssoUrl) == 0 || len(ssoHash) == 0) {
		return _error.NotAuth(ctx)
	}
	oldSsoHash, err := getSSOHash(ctx)
	//如果用户已登录且SSO未发生变更则不进行SSO解析
	if err == nil && len(oldSsoHash) > 0 && loginuser != nil && ssoHash == oldSsoHash {
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

	if loginuser == nil {
		loginuser = new(model.UserInfo)
	}
	uid := ssoUser.Uid
	if ssoUser.Cid > 0 {
		uid = ssoUser.Cid
	}
	loginuser.PID = ssoUser.Pid
	loginuser.UID = uid
	loginuser.Account = ssoUser.Account

	//cid存入cookie
	setSSOUrl(ctx, uid, ssoUrl)
	//用户登录
	contract_expert.SetExpertUserInfo(ctx, loginuser)
	setSSOHash(ctx, ssoHash)
	middleware.Next(ctx)
	return nil
}

func NewExpertSSOMiddleware() *ExpertSSOMiddleware {
	return &ExpertSSOMiddleware{}
}
