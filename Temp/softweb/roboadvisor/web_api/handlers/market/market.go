package market

import (
	"github.com/devfeel/dotweb"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/market"
	"strconv"
)


//Market.GetAtmosphere
func GetAtmosphere(ctx dotweb.Context) error {
	response := contract.NewApiResponse()
	curservice := service.MarketConService()

	objstr,err := curservice.GetAtmosphere()

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	//return ctx.WriteJson(response)
	return ctx.WriteString(objstr)

}


//Market.GetStrategyPool
func GetStrategyPool(ctx dotweb.Context) error {
	response := contract.NewApiResponse()

	StrategyId := ctx.QueryString("strategyid")
	DayCountStr := ctx.QueryString("daycount")
	DayCount,_ := strconv.Atoi(DayCountStr)
	curservice := service.MarketConService()
	objstr,err := curservice.GetStrategyPool(StrategyId,DayCount)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	//return ctx.WriteJson(response)
	return ctx.WriteString(objstr)

}

//Market.BigOrderStrategyHero  英雄榜
func BigOrderStrategyHero(ctx dotweb.Context) error {
	response := contract.NewApiResponse()

	StrategyId := ctx.QueryString("strategyid")
	curservice := service.MarketConService()

	var objstr string
	var err error

	objstr,err = curservice.BigOrderStrategyHero(StrategyId)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	//return ctx.WriteJson(response)
	return ctx.WriteString(objstr)

}

//Market.UnusualGroup
func UnusualGroup(ctx dotweb.Context) error{
	response := contract.NewApiResponse()

	DateStr := ctx.QueryString("date")
	curservice := service.MarketConService()

	objstr,err := curservice.UnusualGroup(DateStr)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	//return ctx.WriteJson(response)
	return ctx.WriteString(objstr)
}


//Market.GetTodayStrategy
func GetTodayStrategy(ctx dotweb.Context) error{
	response := contract.NewApiResponse()

	StrategyId := ctx.QueryString("strategyid")
	curservice := service.MarketConService()
	objstr,err := curservice.GetTodayStrategy(StrategyId)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr

	//return ctx.WriteJson(response)
	return ctx.WriteString(objstr)
}

//Market.GetStrategyList
func GetStrategyList(ctx dotweb.Context) error{
	response := contract.NewApiResponse()

	StrategyId := ctx.QueryString("strategyid")
	curservice := service.MarketConService()
	objstr,err := curservice.GetStrategyList(StrategyId)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	//return ctx.WriteJson(response)
	return ctx.WriteString(objstr)
}


//回踩预选池
func ActiveGoodsPrimary(ctx dotweb.Context) error{

	response := contract.NewApiResponse()

	StrategyId := ctx.QueryString("strategyid")
	curservice := service.MarketConService()
	objstr,err := curservice.ActiveGoodsPrimary(StrategyId)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	return ctx.WriteString(objstr)
}


//传递Key的通用策略股池
func CommonPool(ctx dotweb.Context) error{
	response := contract.NewApiResponse()

	key := ctx.QueryString("key")
	curservice := service.MarketConService()
	objstr,err := curservice.CommonPool(key)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	return ctx.WriteString(objstr)
}


//短线宝
func StrategyOutCommonPool(ctx dotweb.Context) error{
	response := contract.NewApiResponse()

	key := ctx.QueryString("key")
	curservice := service.MarketConService()
	objstr,err := curservice.InnerOutPool(key)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	return ctx.WriteString(objstr)
}


//盘中预警股池 日期
func OnTradingPool(ctx dotweb.Context) error{
	response := contract.NewApiResponse()

	key := ctx.QueryString("key")
	date := ctx.QueryString("date")

	curservice := service.MarketConService()
	objstr,err := curservice.OnTradingMonitorPool(key,date)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	return ctx.WriteString(objstr)
}

//T+1 盘后选股策略股池 日期
func AfterTradingPool(ctx dotweb.Context) error{
	response := contract.NewApiResponse()

	key := ctx.QueryString("key")
	date := ctx.QueryString("date")

	curservice := service.MarketConService()
	objstr,err := curservice.AfterTradingStrategyPool(key,date)

	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = objstr
	return ctx.WriteString(objstr)
}


