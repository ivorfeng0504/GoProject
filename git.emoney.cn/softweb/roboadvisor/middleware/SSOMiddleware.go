package middleware

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/service/mobile"
	"git.emoney.cn/softweb/roboadvisor/protected/service/sso"
	"git.emoney.cn/softweb/roboadvisor/util/crypto"
	"git.emoney.cn/softweb/roboadvisor/util/strings"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/error"
	"github.com/devfeel/dotweb"
	"strings"
)

type SSOMiddleware struct {
	dotweb.BaseMiddlware
}

const (
	RoomListCacheKey  = "UserRoomListCache"
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
	user = nil

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

	nickName := "益盟用户"
	accountLen := _strings.StringLen(account)
	if accountLen > 0 {
		nickName = account
		//如果是QQ号登陆 昵称存储在Pwd字段中
		if strings.Index(nickName, "qq[UID[") > -1 {
			nickName = ssoUser.Pwd
		}
		accountLen = _strings.StringLen(nickName)
		nickNameTmp := nickName
		if accountLen >= 3 {
			mask := "****"
			if accountLen < 7 {
				mask = strings.Repeat("*", accountLen-3)
			}
			nickName = _strings.SubString(nickNameTmp, 0, 3) + mask
			if accountLen > 7 {
				start := accountLen - 3
				if start < 7 {
					start = 7
				}
				nickName += _strings.SubString(nickNameTmp, start, -1)
			}
		}
	}

	//1、首先根据UID尝试登录
	if user == nil && uid > 0 {
		user, err = agent.GetUserByUID(uid)
		if err != nil {
			global.InnerLogger.ErrorFormat(err, "获取用户信息异常 ssoUrl=%s", ssoUrl)
			return _error.ServerError(ctx)
		}
	}

	//2、账号不存在则根据账号登录
	if user == nil && len(account) > 0 {
		//加密并编码成base64
		account, err = mobile.EncryptMobileBase64(account)
		if err != nil {
			return err
		}

		user, err = agent.GetUserByAccount(account)
		if err != nil {
			global.InnerLogger.ErrorFormat(err, "获取用户信息异常 ssoUrl=%s", ssoUrl)
			return _error.ServerError(ctx)
		}
	}

	//自动注册或更新用户信息
	//如果user为空则进行自动注册
	//如果user不为空但是Account为空并且SSO中账号信息不为空 则更新用户的账号信息
	if (user == nil && len(account) > 0 && uid > 0) || (user != nil && len(user.Account) == 0 && len(account) > 0 && uid > 0) {
		userId, err := agent.AddUser(account, nickName, uid)
		if err != nil {
			global.InnerLogger.ErrorFormat(err, "自动注册异常 ssoUrl=%s req=%s", ssoUrl, ctx.Request().Url())
			return _error.ServerError(ctx)
		}
		//再查询一次用户
		user, err = agent.GetUserById(userId)
		if err != nil {
			global.InnerLogger.ErrorFormat(err, "获取用户信息异常2 ssoUrl=%s", ssoUrl)
			return _error.ServerError(ctx)
		}
	}
	if user == nil {
		return _error.NotAuth(ctx)
	}
	//赋值PID
	user.PID = ssoUser.Pid

	//获取用户直播间权限 智投的用户名则为手机号
	roomList, err := ctx.Cache().Get(ctx.SessionID() + RoomListCacheKey)
	if err != nil || roomList == nil {
		roomList, err = agent.GetUserRoomList(user.Account)
		if err != nil {
			global.InnerLogger.Error(err, "获取用户直播间权限异常")
			return _error.ServerError(ctx)
		} else {
			ctx.Cache().Set(ctx.SessionID()+RoomListCacheKey, roomList, 60)
		}
	}
	user.RoomList = roomList.([]int)

	//用户登录
	contract.SetUserInfo(ctx, user)
	setSSOHash(ctx, ssoHash)
	middleware.Next(ctx)
	return nil
}

func NewSSOMiddleware() *SSOMiddleware {
	return &SSOMiddleware{}
}
