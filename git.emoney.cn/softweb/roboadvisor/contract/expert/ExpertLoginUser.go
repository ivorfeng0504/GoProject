package contract_expert

import (
	"github.com/devfeel/dotweb"
	model2 "git.emoney.cn/softweb/roboadvisor/protected/model/expert"
)

const CurrentUserInfoSessionKey = "CurrentExpertUserInfo"

// UserInfo 获取当前登录用户信息
func ExpertUserInfo(ctx dotweb.Context) *model2.UserInfo {
	value := ctx.Session().Get(CurrentUserInfoSessionKey)
	if value == nil {
		return nil
	}
	user := value.(*model2.UserInfo)
	return user
}

// SetUserInfo 设置登录用户信息
func SetExpertUserInfo(ctx dotweb.Context, user *model2.UserInfo) {
	ctx.Session().Set(CurrentUserInfoSessionKey, user)
}

// Logout 退出登录
func Logout(ctx dotweb.Context) {
	ctx.Session().Remove(CurrentUserInfoSessionKey)
}
