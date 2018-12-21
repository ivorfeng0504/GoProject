package middleware

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
)

type UserHomeAuthMiddleware struct {
	dotweb.BaseMiddlware
}

func (middleware *UserHomeAuthMiddleware) Handle(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	if user == nil {
		response.RetCode = -999
		response.RetMsg = "用户未登录"
		return ctx.WriteJson(response)
	}
	middleware.Next(ctx)
	return nil
}

func NewUserHomeAuthMiddleware() *UserHomeAuthMiddleware {
	return &UserHomeAuthMiddleware{}
}
