package userhome

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"time"
)

// InsertOrUpdateUserStock 用户领取奖品
func InsertOrUpdateUserStock(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	requestData := agent.InsertUserStockRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if requestData.UserInfo == nil {
		err = errors.New("未获取到用户信息")
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if requestData.ActivityId <= 0 {
		err = errors.New("未获取到活动Id")
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := service.NewUserReceivedStockService()
	//查询用户当天是否已经领取了研报
	userStock, err := srv.GetUserStockToday(requestData.UserInfo.ID)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	//如果已经领取了研报并且研报内容没有发生变更则直接返回，否则更新用户领取的研报内容
	//若没有则获取当日研报并写入数据库
	receiveStockActivitySrv := service.NewReceiveStockActivityService()
	stockPool, err := receiveStockActivitySrv.GetNewstStockPool()
	if err != nil {
		response.RetCode = -6
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if userStock == nil {
		//新增
		userStock = &userhome_model.UserReceivedStock{
			UID:                    requestData.UserInfo.UID,
			UserInfoId:             requestData.UserInfo.ID,
			ReceiveStockActivityId: stockPool.ReceiveStockActivityId,
			StockCreateTime:        stockPool.CreateTime,
			StockList:              stockPool.StockList,
			IssueNumber:            time.Now().Format("20060102"),
			ReportUrl:              stockPool.ReportUrl,
		}
		_, err = srv.InsertUserStock(userStock, requestData.ActivityId)
		if err != nil {
			response.RetCode = -7
			response.RetMsg = err.Error()
			return ctx.WriteJson(response)
		}
	} else {
		//有变动则更新
		if len(stockPool.ReportUrl) > 0 && stockPool.ReportUrl != userStock.ReportUrl {
			err = srv.UpdateUserReceivedStock(requestData.UserInfo.ID, userStock.UserReceivedStockId, stockPool.ReportUrl)
			if err != nil {
				response.RetCode = -8
				response.RetMsg = err.Error()
				return ctx.WriteJson(response)
			} else {
				userStock.ReportUrl = stockPool.ReportUrl
			}
		}
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = userStock
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetUserStockToday 获取用户当日领取的股票
func GetUserStockToday(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	userInfoIdF, success := request.RequestData.(float64)
	if !success {
		response.RetCode = -2
		response.RetMsg = "用户Id不正确"
		return ctx.WriteJson(response)
	}
	userInfoId := int(userInfoIdF)
	srv := service.NewUserReceivedStockService()
	userStock, err := srv.GetUserStockToday(userInfoId)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = userStock
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetUserStockHistory 获取用户领取的股票历史
func GetUserStockHistory(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := agent.UserReceivedStockHistoryRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := service.NewUserReceivedStockService()
	userStockList, err := srv.GetUserStockHistory(requestData.UserInfoId, requestData.Top)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = userStockList
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
