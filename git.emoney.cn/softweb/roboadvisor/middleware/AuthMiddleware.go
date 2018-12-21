package middleware

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"github.com/devfeel/dotweb"
)

type AuthMiddleware struct {
	dotweb.BaseMiddlware
}

func (middleware *AuthMiddleware) Handle(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract.UserInfo(ctx)
	if user == nil {
		response.RetCode = -999
		response.RetMsg = "用户未登录"
		return ctx.WriteJson(response)
	}
	middleware.Next(ctx)
	return nil
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}
