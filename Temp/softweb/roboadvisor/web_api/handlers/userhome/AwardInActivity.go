package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// GetActivityAwardListByActivityId 获取指定活动的奖品列表
func GetActivityAwardListByActivityId(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := agent.ActivityAwardRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := service.NewAwardInActivityService()
	awardList, err := srv.GetActivityAwardListByActivityId(requestData.ActivityId, requestData.IgnoreExpired)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		response.Message = false
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = awardList
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
