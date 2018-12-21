package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userloginlog"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// GetUserSerialLoginRule 获取用户连登奖励
func GetUserSerialLoginRule(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := userhome_model.UserInfo{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := service.NewSerialLoginRuleService()
	rule, err := srv.GetUserSerialLoginRule(requestData)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		response.Message = false
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = rule
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetUserLevel 获取用户等级
func GetUserLevel(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := userhome_model.UserInfo{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	loginCount, err := userloginlog.GetUserLoginCount(requestData.GID)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	srv := service.NewSerialLoginRuleService()
	level := srv.GetUserLevelByLoginDay(loginCount)

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = level

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
