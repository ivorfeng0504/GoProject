package middleware

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/myoptional"
	"github.com/devfeel/dotweb"
)

type MyOptionalAuthMiddleware struct {
	dotweb.BaseMiddlware
}

func (middleware *MyOptionalAuthMiddleware) Handle(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_myoptional.UserInfo(ctx)
	if user == nil {
		response.RetCode = -999
		response.RetMsg = "用户未登录"
		return ctx.WriteJson(response)
	}
	middleware.Next(ctx)
	return nil
}

func NewMyOptionalAuthMiddleware() *MyOptionalAuthMiddleware {
	return &MyOptionalAuthMiddleware{}
}
