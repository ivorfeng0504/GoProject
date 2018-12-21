package contract_train

import (
	"github.com/devfeel/dotweb"
	livemodel "git.emoney.cn/softweb/roboadvisor/protected/model/live"
)

const CurrentUserInfoSessionKey = "CurrentTrainUserInfo"

// UserInfo 获取当前登录用户信息
func TrainUserInfo(ctx dotweb.Context) *livemodel.User {
	value := ctx.Session().Get(CurrentUserInfoSessionKey)
	if value == nil {
		return nil
	}
	user := value.(*livemodel.User)
	return user
}

// SetUserInfo 设置登录用户信息
func SetTrainUserInfo(ctx dotweb.Context, user *livemodel.User) {
	ctx.Session().Set(CurrentUserInfoSessionKey, user)
}

// Logout 退出登录
func Logout(ctx dotweb.Context) {
	ctx.Session().Remove(CurrentUserInfoSessionKey)
}
