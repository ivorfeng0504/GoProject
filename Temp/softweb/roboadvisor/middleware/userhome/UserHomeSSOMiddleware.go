package middleware

import (
	"fmt"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/mobile"
	"git.emoney.cn/softweb/roboadvisor/protected/service/sso"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/crypto"
	"git.emoney.cn/softweb/roboadvisor/web_front/handlers/error"
	"git.emoney.cn/softweb/roboadvisor/web_front_userhome/validate"
	"github.com/devfeel/dotweb"
	"net/url"
	"strconv"
	"strings"
)

type UserHomeSSOMiddleware struct {
	dotweb.BaseMiddlware
}

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
	token := url.QueryEscape(ctx.QueryString("token"))
	bata := ctx.QueryString("bata")
	url := fmt.Sprintf(ssoTemplate, Version, uid, cid, pid, sid, tid, agentid, clienttype, OutOfDate, pluglet, token, bata)
	ssoHash = _crypto.MD5(url)[:8]
	return url, ssoHash
}

func getSSOHash(ctx dotweb.Context) (ssoHash string, err error) {
	ssoHash, err = ctx.ReadCookieValue("userhomessohash")
	return ssoHash, err
}

func setSSOHash(ctx dotweb.Context, ssoHash string, ssoUrl string) {
	ctx.SetCookieValue("userhomessohash", ssoHash, 0)
	//fmt.Println(ssoUrl)
	ctx.SetCookieValue("userhome.ssourl", ssoUrl, 0)
}

func (middleware *UserHomeSSOMiddleware) Handle(ctx dotweb.Context) error {
	ssoUrl, ssoHash := getSSOUrl(ctx)

	loginuser := contract_userhome.UserHomeUserInfo(ctx)
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
	//如果用户已登录无sso信息则直接读取session数据
	if err == nil && len(oldSsoHash) > 0 && loginuser != nil && ctx.QueryString("token") == "" {
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
	TIDStr := ctx.QueryString("tid")
	tid, err := strconv.Atoi(TIDStr)
	if err != nil {
		tid = 0
	}
	SIDStr := ctx.QueryString("sid")
	sid, err := strconv.Atoi(SIDStr)
	if err != nil {
		sid = 0
	}

	uid := ssoUser.Uid
	gid := ssoUser.Uid
	if ssoUser.Cid > 0 {
		uid = ssoUser.Cid
	}
	account := ssoUser.Account
	username := ssoUser.Account
	pid := ssoUser.Pid
	mobilex := ""
	mobilemask := ""
	nickName := ""
	usertype := _const.UserHome_UserType_Free

	if len(account) == 0 && uid <= 0 {
		global.InnerLogger.DebugFormat("SSO 解析异常  Account与UID同时不存在 ssoUrl=%s", ssoUrl)
		return _error.NotAuth(ctx)
	}

	//QQ号单独处理- [替换为_
	if strings.Index(username,"qq[UID[")>-1 && len(username)>30 {
		username = strings.Replace(username, "[", "_", -1)
		account = strings.Replace(username, "[", "_", -1)
	}

	service := service.NewUserInfoService()

	//判断用户类型（游客、第三方（微信 QQ）、手机）
	if strings.Contains(_const.ZTPID_YK,strconv.Itoa(pid)) {
		nickName = "游客"
		usertype = _const.UserHome_UserType_Free
	} else {
		if len(account) > 0 {
			if len(account) >= 3 {
				nickName = account[:3] + "****"
				if len(account) > 7 {
					nickName += account[7:]
				}
			} else {
				nickName = account
			}
		}

		if validate.IsMobile(account) {
			//手机号-加密账号
			usertype = _const.UserHome_UserType_Mobile
			mobilemask = nickName
			//加密并编码成base64
			mobilex, err = mobile.EncryptMobileHex(account)
			if err != nil {
				return err
			}
			account = mobilex

		} else {
			usertype = _const.UserHome_UserType_EM

			//QQ和微信登录SSO密码一栏 存入的是昵称（此处昵称特殊取值）
			if strings.Contains(account, "wx_") {
				usertype = _const.UserHome_UserType_WeChat
				nickName = ssoUser.Pwd
			}
			if strings.Contains(account, "qq_") {
				usertype = _const.UserHome_UserType_QQ
				nickName = ssoUser.Pwd
			}
		}
	}
	//1、根据uid登录
	if uid > 0 {
		loginuser, _ = service.GetUserInfoByUID(uid)
	}
	//自动注册或更新用户信息
	//如果user为空则进行自动注册
	if loginuser == nil || loginuser.UID <= 0 {
		newuser := new(model.UserInfo)
		newuser.UID = uid
		newuser.NickName = nickName
		newuser.Account = account

		if validate.IsMobile(account) {
			newuser.Account = mobilex
		}
		//手机账号
		if len(mobilex) > 0 {
			newuser.MobileX = mobilex
			newuser.MobileMask = mobilemask
		}
		//微信账号
		if usertype == _const.UserHome_UserType_WeChat {
			newuser.OpenID_WeChat = account
		}
		if usertype == _const.UserHome_UserType_QQ {
			newuser.OpenID_QQ = account
		}
		newuser.UserType = usertype

		_, err := service.AddUserInfo(newuser)
		if err != nil {
			global.InnerLogger.Error(err, "自动注册异常")
			return _error.ServerError(ctx)
		}

	} else {
		//更新最后登录时间
		_, err := service.ModifyLastLoginTime(uid)
		if err != nil {
			global.InnerLogger.Error(err, "更新最后登录时间异常")
			return _error.ServerError(ctx)
		}
	}

	//再查询一次用户
	loginuser, _ = service.GetUserInfoByUID(uid)

	if loginuser == nil {
		return _error.NotAuth(ctx)
	}

	//用户类型重新赋值
	loginuser.UserType = usertype

	//判断用户所属版本（初付费和体验期，其他都归属免费版）
	loginuser.PIDType = _const.ZTUserType_Free
	if strings.Contains(_const.ZTPID_Pay, strconv.Itoa(pid)) { //付费
		loginuser.PIDType = _const.ZTUserType_Pay
	}
	if strings.Contains(_const.ZTPID_Experience, strconv.Itoa(pid)) { //体验期内
		loginuser.PIDType = _const.ZTUserType_Experience
	}

	loginuser.PID = pid
	loginuser.SID = sid
	loginuser.TID = tid
	loginuser.GID = gid
	loginuser.UserName = username
	//用户登录
	contract_userhome.SetUserHomeUserInfo(ctx, loginuser)
	setSSOHash(ctx, ssoHash, ssoUrl)
	middleware.Next(ctx)
	return nil
}

func NewUserHomeSSOMiddleware() *UserHomeSSOMiddleware {
	return &UserHomeSSOMiddleware{}
}
