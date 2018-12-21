package stocknews

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	myoptional_srv "git.emoney.cn/softweb/roboadvisor/protected/service/myoptional"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// GetStockNewsInfo 获取股票相关资讯
func GetStockNewsInfo(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StockNewsInformationRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := myoptional_srv.NewStockNewsInformationService()
	newsList, err := srv.GetStockNewsInfo(requestData.StartTime, requestData.StockList, requestData.Top)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetStockNewsInfoByHSet 根据股票代码获取相关资讯（直接读HSet缓存）
func GetStockNewsInfoByHSet(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData agent.StockNewsInformationRequest
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := myoptional_srv.NewStockNewsInformationService()
	newsList, err := srv.GetStockNewsInfoByHSet(requestData.StockList, requestData.PerTop, requestData.Top)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = newsList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
