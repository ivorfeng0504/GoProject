package stocknews

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/util/http"
	"github.com/devfeel/dotweb"
	"strings"
)

// GetStockNewsInfo 获取股票相关资讯
func GetStockNewsInfo(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	stockListStr := ctx.PostFormValue("StockList")
	if len(stockListStr) == 0 {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		return ctx.WriteJson(response)
	}
	stockList := strings.Split(stockListStr, ",")
	if len(stockList) == 0 {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		return ctx.WriteJson(response)
	}
	result, err := agent.GetStockNewsInfo(stockList)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = result
	}
	return ctx.WriteJson(response)
}

// GetStockNewsInfoByHSet 根据股票代码获取相关资讯（直接读HSet缓存）
func GetStockNewsInfoByHSet(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	stockListStr := _http.GetRequestValue(ctx, "StockList")
	if len(stockListStr) == 0 {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		return ctx.WriteJson(response)
	}
	stockList := strings.Split(stockListStr, ",")
	if len(stockList) == 0 {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		return ctx.WriteJson(response)
	}
	result, err := agent.GetStockNewsInfoByHSet(stockList)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = result
	}
	return ctx.WriteJson(response)
}
