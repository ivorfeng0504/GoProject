package middleware

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract/expert"
)

type ExpertAuthMiddleware struct {
	dotweb.BaseMiddlware
}

func (middleware *ExpertAuthMiddleware) Handle(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_expert.ExpertUserInfo(ctx)
	if user == nil {
		response.RetCode = -999
		response.RetMsg = "用户未登录"
		return ctx.WriteJson(response)
	}
	middleware.Next(ctx)
	return nil
}

func NewExpertAuthMiddleware() *ExpertAuthMiddleware {
	return &ExpertAuthMiddleware{}
}
