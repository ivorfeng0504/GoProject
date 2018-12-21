package contract

import (
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
	"github.com/devfeel/dotweb"
)

const CurrentUserInfoSessionKey = "CurrentUserInfo"

// UserInfo 获取当前登录用户信息
func UserInfo(ctx dotweb.Context) *livemodel.User {
	value := ctx.Session().Get(CurrentUserInfoSessionKey)
	if value == nil {
		return nil
	}
	user := value.(*livemodel.User)
	return user
}

// SetUserInfo 设置登录用户信息
func SetUserInfo(ctx dotweb.Context, user *livemodel.User) {
	ctx.Session().Set(CurrentUserInfoSessionKey, user)
}

// Logout 退出登录
func Logout(ctx dotweb.Context) {
	ctx.Session().Remove(CurrentUserInfoSessionKey)
}
