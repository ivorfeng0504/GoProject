package strategy

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/strategy"
	"strconv"
	"fmt"
)

// 策略详情页
func Strategy(ctx dotweb.Context) error {
	return contract.RenderHtml(ctx, "strategy.html")
}


// 策略详情接口 Strategy Detail
func StrategyDetail(ctx dotweb.Context) error{
	response := contract.NewResonseInfo()

	StrategyIDStr := ctx.QueryString("StrategyID")
	StrategyID, err := strconv.Atoi(StrategyIDStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略ID不正确"
		return ctx.WriteJson(response)
	}

	newsService := service.StrategyInfoService()

	strategy, err := newsService.GetStrategyInfoById(StrategyID)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategy
	return ctx.WriteJson(response)
}



//获取策略信息 软件接口使用
func GetStrategyInfo(ctx dotweb.Context) error {
	response := contract.NewApiResponse()

	StrategyIdStr := ctx.QueryString("StrategyId")
	StrategyId, err := strconv.Atoi(StrategyIdStr)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = "策略编号不正确"
		return ctx.WriteJson(response)
	}

	strategyService := service.StrategyInfoService()
	strategyobj, err := strategyService.GetStrategyInfoById(StrategyId)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyobj
	var o = strategyobj
	return ctx.WriteString(fmt.Sprintf("{|*%d*|}{|*%s*|}{|*%s*|}{|*%s*|}{|*%s*|}{|*%s*|}",o.StrategyID,o.StrategyName,o.LogoUrl,o.Summary,o.Description,o.ImagesUrl))
	//return ctx.WriteJson(response)
}