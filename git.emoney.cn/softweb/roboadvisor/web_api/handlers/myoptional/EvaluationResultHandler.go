package myoptional

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	myoptional_model "git.emoney.cn/softweb/roboadvisor/protected/model/myoptional"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"
	myoptional_srv "git.emoney.cn/softweb/roboadvisor/protected/service/myoptional"
	"git.emoney.cn/softweb/roboadvisor/protected/service/strategyapi"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// GetResult 查看当前用户的评测结果
func GetResult(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	uid, success := request.RequestData.(string)
	if success == false {
		response.RetCode = -2
		response.RetMsg = "获取UID失败"
		return ctx.WriteJson(response)
	}
	evaSrv := myoptional_srv.NewEvaluationResultService()
	result, err := evaSrv.GetEvaluationResult(uid)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = result
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// SubmitResult 提交评测
func SubmitResult(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := myoptional_model.EvaluationResult{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	evaSrv := myoptional_srv.NewEvaluationResultService()
	result, err := evaSrv.GetEvaluationResult(requestData.UID)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if result == nil {
		err = evaSrv.InsertEvaluationResult(requestData)
	} else {
		err = evaSrv.UpdateEvaluationResult(requestData)
	}
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	return ctx.WriteJson(response)
}

// GetStrategyList 获取推荐策略及股票列表
func GetStrategyList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var requestData []*agent.StrategyInfo
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	var strategyInfoList []*agent.StrategyInfo
	for _, strategy := range requestData {
		strategyInfo := &agent.StrategyInfo{
			StrategyId:   strategy.StrategyId,
			StrategyName: strategy.StrategyName,
		}
		if strategy.StrategyKeyList != nil && len(strategy.StrategyKeyList) > 0 {
			for _, strategyKey := range strategy.StrategyKeyList {
				stockList, err := strategyapi.GetCommonPoolStockList(strategy.Top, strategyKey, strategy.ReadCacheMode)
				if err == nil {
					strategyInfo.StockList = append(strategyInfo.StockList, stockList...)
				}
			}
		} else {
			topicSrv := expertnews.NewExpertNews_TopicService()
			topicStockList, err := topicSrv.StatTopicFocusStock()
			if err != nil {
				response.RetCode = -3
				response.RetMsg = err.Error()
				return ctx.WriteJson(response)
			}
			if topicStockList != nil {
				for _, topicStock := range *topicStockList {
					strategyInfo.StockList = append(strategyInfo.StockList, &userhome_model.StockInfo{
						StockCode: topicStock.StockCode,
						StockName: topicStock.StockName,
					})
				}
			}
		}
		strategyInfoList = append(strategyInfoList, strategyInfo)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = strategyInfoList
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
