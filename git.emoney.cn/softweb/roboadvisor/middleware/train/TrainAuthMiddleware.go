package train

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract/train"
)

type TrainAuthMiddleware struct {
	dotweb.BaseMiddlware
}

func (middleware *TrainAuthMiddleware) Handle(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_train.TrainUserInfo(ctx)
	if user == nil {
		response.RetCode = -999
		response.RetMsg = "用户未登录"
		return ctx.WriteJson(response)
	}
	middleware.Next(ctx)
	return nil
}

func NewTrainAuthMiddleware() *TrainAuthMiddleware {
	return &TrainAuthMiddleware{}
}
