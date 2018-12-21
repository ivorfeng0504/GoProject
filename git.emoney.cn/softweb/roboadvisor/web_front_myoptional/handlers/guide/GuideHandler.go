package guide

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"github.com/devfeel/dotweb"
	"time"
)

// GetRecommendStockList 获取推荐策略及股票列表
func GetRecommendStockList(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	callback := ctx.QueryString("callback")
	if len(callback) == 0 {
		response.RetCode = -1
		response.RetMsg = "回调函数不正确，请稍后再试"
		return ctx.WriteJsonp(callback, response)
	}
	strategyList, err := agent.GetRecommendStockList()
	if err != nil {
		response.RetCode = -2
		response.RetMsg = "获取推荐策略及股票列表异常，请稍后再试！"
		return ctx.WriteJsonp(callback, response)
	} else {
		if strategyList != nil {
			for _, strategy := range strategyList {
				if strategy.StockList == nil || len(strategy.StockList) == 0 {
					continue
				}
				//如果股票的创建时间早于3天 则忽略
				var stockListDist []*userhome_model.StockInfo
				for _, stockInfo := range strategy.StockList {
					if time.Now().AddDate(0, 0, -3).Before(stockInfo.CreateTime) {
						stockListDist = append(stockListDist, stockInfo)
					}
				}
				//股票取前10个
				strategy.StockList = stockListDist
				if strategy.StockList != nil && len(strategy.StockList) > 10 {
					strategy.StockList = strategy.StockList[:10]
				}
			}
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyList
	return ctx.WriteJsonp(callback, response)
}
