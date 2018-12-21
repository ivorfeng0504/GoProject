package activity

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/contract/userhome"
	"git.emoney.cn/softweb/roboadvisor/global"
	"git.emoney.cn/softweb/roboadvisor/protected/viewmodel/userhome"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/mapper"
)

// GuessMarket 领数字活动首页
func GuessMarket(ctx dotweb.Context) error {
	return contract_userhome.RenderUserHomeHtml(ctx, "guess_market.html")
}

// GetAwardInfo 获取奖品信息-本周大奖
func GetAwardInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	awardList, err := agent.GetGuessMarketAward()
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取历史信息异常"
		global.InnerLogger.ErrorFormat(err, "获取猜数字奖品信息异常")
		return ctx.WriteJson(response)
	}

	var awardSimpleList []*viewmodel.ActivityAwardSimple
	if awardList != nil && len(awardList) > 0 {
		mapper.MapperSlice(awardList, &awardSimpleList)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = awardSimpleList
	return ctx.WriteJson(response)
}

// GetGuessMarketHistoryInfo 获取历史信息-包括我的抽奖记录&往期中奖号码&中奖名单
func GetGuessMarketHistoryInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	history, err := agent.GetGuessMarketHistoryInfo(user.ID)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取历史信息异常"
		global.InnerLogger.ErrorFormat(err, "获取历史信息异常 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = history
	return ctx.WriteJson(response)
}

// GetGuessMarketNumber 获取抽奖号码
func GetGuessMarketNumber(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	number, err := agent.GetGuessMarketNumber(*user)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取抽奖号码异常，请稍后再试！"
		global.InnerLogger.ErrorFormat(err, "获取抽奖号码异常 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = number
	return ctx.WriteJson(response)
}

// GetGuessMarketNumberWithDetail 获取抽奖号码并返回当前活动的详细信息
func GetGuessMarketNumberWithDetail(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	number, err := agent.GetGuessMarketNumber(*user)
	if err != nil || len(number) != 4 {
		response.RetCode = -1
		response.RetMsg = "获取抽奖号码异常，请稍后再试！"
		global.InnerLogger.ErrorFormat(err, "获取抽奖号码异常 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}

	currentGuessInfo, err := agent.GetCurrentGuessInfo(user.ID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取当前的抽奖信息异常，请稍后再试！"
		global.InnerLogger.ErrorFormat(err, "获取当前的抽奖信息异常 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = currentGuessInfo
	return ctx.WriteJson(response)
}

// GetCurrentGuessInfo 获取当前的抽奖信息-抽奖号码，开奖时间等信息
func GetCurrentGuessInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	user := contract_userhome.UserHomeUserInfo(ctx)
	currentGuessInfo, err := agent.GetCurrentGuessInfo(user.ID)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "获取当前的抽奖信息异常，请稍后再试！"
		global.InnerLogger.ErrorFormat(err, "获取当前的抽奖信息异常 userInfoId=%d", user.ID)
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = currentGuessInfo
	return ctx.WriteJson(response)
}
