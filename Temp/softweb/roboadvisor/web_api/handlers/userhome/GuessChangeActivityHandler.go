package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
)

// GetCurrentGuessChange 获取当期猜涨跌活动情况
func GetCurrentGuessChange(ctx dotweb.Context) error {
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
	if requestData.ID <= 0 {
		response.RetCode = -3
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}

	currentGuessChageInfo := agent.CurrentGuessChange{}
	//查询当期活动信息
	guessChangeActivitySrv := service.NewGuessChangeActivityService()
	activity, err := guessChangeActivitySrv.GetCurrentGuessChangeActivity()
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if activity == nil {
		response.RetCode = -5
		response.RetMsg = "当前没有猜涨跌活动"
		return ctx.WriteJson(response)
	}
	//填充活动信息
	currentGuessChageInfo.IssueNumber = activity.IssueNumber

	//查询当期活动参与统计
	userInGuessSrv := service.NewUserInGuessChangeActivityService()
	total, err := userInGuessSrv.GetGuessTotal(currentGuessChageInfo.IssueNumber)
	if err != nil {
		response.RetCode = -6
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if total == nil {
		response.RetCode = -7
		response.RetMsg = "未获取到当前活动的统计数据"
		return ctx.WriteJson(response)
	}
	currentGuessChageInfo.GuessUp, currentGuessChageInfo.GuessDown = total.UpCount, total.DownCount

	//查询当前我的参与信息
	myJoinRecord, err := userInGuessSrv.GetUserInGuessChangeActivity(requestData.ID, currentGuessChageInfo.IssueNumber)
	if err != nil {
		response.RetCode = -8
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if myJoinRecord != nil {
		currentGuessChageInfo.MyGuessResult = myJoinRecord.Result
		currentGuessChageInfo.IsJoin = true
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = currentGuessChageInfo

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GuessChangeSubmit 提交猜涨跌竞猜结果
func GuessChangeSubmit(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := agent.GuessChangeSubmitRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if requestData.UserInfo.ID <= 0 {
		response.RetCode = -3
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}
	currentGuessChageInfo := agent.CurrentGuessChange{}
	//查询当期活动信息
	guessChangeActivitySrv := service.NewGuessChangeActivityService()
	activity, err := guessChangeActivitySrv.GetCurrentGuessChangeActivity()
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if activity == nil {
		response.RetCode = -5
		response.RetMsg = "当前没有猜涨跌活动"
		return ctx.WriteJson(response)
	}
	//填充活动信息
	currentGuessChageInfo.IssueNumber = activity.IssueNumber

	//参与竞猜
	userInGuessSrv := service.NewUserInGuessChangeActivityService()
	err = userInGuessSrv.GuessChangeSubmit(currentGuessChageInfo.IssueNumber, requestData.UserInfo.ID, requestData.UserInfo.UID, requestData.UserInfo.NickName, requestData.Result, activity.ActivityCycle, requestData.ActivityId)
	if err != nil {
		response.RetCode = -6
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	currentGuessChageInfo.MyGuessResult = requestData.Result
	currentGuessChageInfo.IsJoin = true

	//查询当期活动参与统计
	total, err := userInGuessSrv.GetGuessTotal(currentGuessChageInfo.IssueNumber)
	if err != nil {
		response.RetCode = -7
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if total == nil {
		response.RetCode = -8
		response.RetMsg = "未获取到当前活动的统计数据"
		return ctx.WriteJson(response)
	}
	currentGuessChageInfo.GuessUp, currentGuessChageInfo.GuessDown = total.UpCount, total.DownCount

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = currentGuessChageInfo

	return ctx.WriteJson(response)
}

// GetMyGuessChangeInfoCurrentWeek 获取本周我的竞猜记录
func GetMyGuessChangeInfoCurrentWeek(ctx dotweb.Context) error {
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
	if requestData.ID <= 0 {
		response.RetCode = -3
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}

	userInGuessChangeSrv := service.NewUserInGuessChangeActivityService()
	myGuessChangeList, err := userInGuessChangeSrv.GetMyGuessChangeInfoCurrentWeek(requestData.ID)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = myGuessChangeList

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetMyGuessageAward 获取我的猜涨跌奖品
func GetMyGuessageAwardList(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := agent.MyGuessChangeInfoRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if requestData.UserInfo.ID <= 0 {
		response.RetCode = -3
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}

	userInGuessChangeSrv := service.NewUserInGuessChangeActivityService()
	myGuessChangeList, err := userInGuessChangeSrv.GetMyGuessageAwardList(requestData.UserInfo.ID, requestData.Top)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = myGuessChangeList

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetMyGuessChangeInfoNewst 获取我最近的猜涨跌记录
func GetMyGuessChangeInfoNewst(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := agent.MyGuessChangeInfoRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	if requestData.UserInfo.ID <= 0 {
		response.RetCode = -3
		response.RetMsg = "未获取到用户信息"
		return ctx.WriteJson(response)
	}

	userInGuessChangeSrv := service.NewUserInGuessChangeActivityService()
	myGuessChangeList, err := userInGuessChangeSrv.GetMyGuessChangeInfoNewst(requestData.UserInfo.ID, requestData.Top)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = myGuessChangeList

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
