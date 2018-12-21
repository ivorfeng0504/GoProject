package contract

import (
	"github.com/devfeel/dotweb"
	"emoney.cn/fundchannel/protected/model"
)

const (
	CurrentUserInfoSessionKey = "CurrentUserInfo"
	CurrentSSOQuerySessionKey = "CurrentSSOQuerySessionKey"
)

// UserInfo 获取当前登录用户信息
func UserInfo(ctx dotweb.Context) *fund.User {
	value := ctx.Session().Get(CurrentUserInfoSessionKey)
	if value == nil {
		return nil
	}
	user := value.(*fund.User)
	return user
}

// SetUserInfo 设置登录用户信息
func SetUserInfo(ctx dotweb.Context, user *fund.User) {
	ctx.Session().Set(CurrentUserInfoSessionKey, user)
}

// SetSSOQuery 记录登录用户SSO串
func SetSSOQuery(ctx dotweb.Context, ssoQuery string) {
	ctx.Session().Set(CurrentSSOQuerySessionKey, ssoQuery)
}
func GetSSOQuery(ctx dotweb.Context) (ssoQuery string) {
	return ctx.Session().GetString(CurrentSSOQuerySessionKey)
}

// Logout 退出登录
func Logout(ctx dotweb.Context) {
	ctx.Session().Remove(CurrentUserInfoSessionKey)
	ctx.Session().Remove(CurrentSSOQuerySessionKey)
}
