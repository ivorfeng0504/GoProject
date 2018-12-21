package contract_userhome

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
)

const CurrentUserInfoSessionKey = "CurrentUserHomeUserInfo"

// UserInfo 获取当前登录用户信息
func UserHomeUserInfo(ctx dotweb.Context) *model.UserInfo {
	value := ctx.Session().Get(CurrentUserInfoSessionKey)
	if value == nil {
		return nil
	}
	user := value.(*model.UserInfo)
	return user
}

// SetUserInfo 设置登录用户信息
func SetUserHomeUserInfo(ctx dotweb.Context, user *model.UserInfo) {
	ctx.Session().Set(CurrentUserInfoSessionKey, user)
}

// Logout 退出登录
func Logout(ctx dotweb.Context) {
	ctx.Session().Remove(CurrentUserInfoSessionKey)
}
