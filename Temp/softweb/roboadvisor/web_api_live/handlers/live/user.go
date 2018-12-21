package live

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	liveservice "git.emoney.cn/softweb/roboadvisor/protected/service/live"
	"git.emoney.cn/softweb/roboadvisor/protected/service/mobile"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// GetUserById 通过用户Id获取用户信息
func GetUserById(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	userIdF, success := request.RequestData.(float64)
	if !success {
		response.RetCode = -2
		response.RetMsg = "UserId不正确"
		return ctx.WriteJson(response)
	}
	userId := int(userIdF)
	if userId <= 0 {
		response.RetCode = -2
		response.RetMsg = "UserId不正确"
		return ctx.WriteJson(response)
	}

	userService := liveservice.NewUserService()
	user, err := userService.GetUserById(userId)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = user
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetUserByUID 通过UID账号获取用户信息
func GetUserByUID(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	uidF, success := request.RequestData.(float64)
	if !success {
		response.RetCode = -2
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}
	uid := int64(uidF)
	if uid <= 0 {
		response.RetCode = -2
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}

	userService := liveservice.NewUserService()
	user, err := userService.GetUserByUID(uid)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = user
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetUserByAccount 通过账号获取用户信息
func GetUserByAccount(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	account := request.RequestData.(string)
	if len(account) == 0 {
		response.RetCode = -2
		response.RetMsg = "Account不正确"
		return ctx.WriteJson(response)
	}

	account = mobile.Hex2Base64(account)
	userService := liveservice.NewUserService()
	user, err := userService.GetUserByAccount(account)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = user
	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// AddUser 通过账号和UID自动注册用户
func AddUser(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	requestData := agent.AddUserRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	if len(requestData.Account) == 0 {
		response.RetCode = -2
		response.RetMsg = "Account不正确"
		return ctx.WriteJson(response)
	}

	if requestData.UID <= 0 {
		response.RetCode = -3
		response.RetMsg = "UID不正确"
		return ctx.WriteJson(response)
	}

	if len(requestData.Source) < 5 {
		response.RetCode = -4
		response.RetMsg = "用户来源不正确"
		return ctx.WriteJson(response)
	}

	userService := liveservice.NewUserService()
	//十六进制编码成base64
	requestData.Account = mobile.Hex2Base64(requestData.Account)

	userId, err := userService.AddUser(requestData.Account, requestData.NickName, requestData.UID, requestData.Source)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = userId
	return ctx.WriteJson(response)
}
