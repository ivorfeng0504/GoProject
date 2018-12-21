package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// GetActivityList 获取活动列表
func GetActivityList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	activityStateF, success := request.RequestData.(float64)
	if !success {
		response.RetCode = -2
		response.RetMsg = "活动状态不正确"
		return ctx.WriteJson(response)
	}
	activityState := int(activityStateF)
	activitySrv := service.NewActivityService()
	activityList, err := activitySrv.GetActivityList(activityState)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		response.Message = false
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = activityList
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetUserMedal 获取用户勋章
func GetUserMedal(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := model.UserInfo{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	userMedalSrv := service.NewUserMedalService()
	medalInfo := userMedalSrv.GetUserMedalInfoAll(requestData)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		response.Message = false
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = medalInfo
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
