package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"github.com/devfeel/dotweb"
)

// GetUserInActivityMap 获取用户参与活动Map
func GetUserInActivityMap(ctx dotweb.Context) error {
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
	userInActSrv := service.NewUserInActivityService()
	userActivityMap, err := userInActSrv.GetUserInActivityMap(userInfoId)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		response.Message = false
	} else {
		response.RetCode = 0
		response.RetMsg = "SUCCESS"
		response.Message = userActivityMap
	}

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
