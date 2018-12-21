package myoptional

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	expertnews2 "git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	"github.com/devfeel/dotweb"
	"strconv"
	"git.emoney.cn/softweb/roboadvisor/contract/myoptional"
)

func GetFocusStrategyList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_myoptional.UserInfo(ctx)

	UIDStr := ctx.QueryString("UID")
	UID, err := strconv.ParseInt(UIDStr, 10, 64)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}

	if user != nil {
		UID = user.UID
	}

	strategyService := expertnews2.NewExpertNews_FocusStrategyService()
	focusStrategyList, hasFocus, err := strategyService.GetFocusStrategyByUID(UID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = strconv.FormatBool(hasFocus)
	response.Message = focusStrategyList
	return ctx.WriteJson(response)
}

