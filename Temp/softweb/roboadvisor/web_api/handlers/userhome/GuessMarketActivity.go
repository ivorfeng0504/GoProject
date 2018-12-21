package userhome

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/const"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"git.emoney.cn/softweb/roboadvisor/protected/service/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
	"github.com/devfeel/dotweb"
	"time"
)

// GetGuessMarketHistoryInfo 获取历史信息-包括我的抽奖记录&往期中奖号码&中奖名单
func GetGuessMarketHistoryInfo(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	history := agent.GuessMarketActivityHistory{}
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := agent.GuessMarketActivityHistoryRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	guessMarketActivitySrv := service.NewGuessMarketActivityService()
	userInGuessMarketActivitySrv := service.NewUserInGuessMarketActivityService()
	history.HistoryLuckNumberList, err = guessMarketActivitySrv.GetGuessMarketActivityList(requestData.HistoryLuckNumberCount)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	history.MyGuessCount, err = userInGuessMarketActivitySrv.GetUserJoinCount(requestData.UserInfoId)
	if err != nil {
		response.RetCode = -4
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	history.MyGuessList, err = userInGuessMarketActivitySrv.GetUserInGuessMarketActivityList(requestData.MyGuessCount, requestData.UserInfoId)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	history.WinnerHistoryList, err = userInGuessMarketActivitySrv.GetUserInGuessMarketActivityGuessedList(requestData.WinnerHistoryCount)
	if err != nil {
		response.RetCode = -6
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = history

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetGuessMarketNumber 用户获取抽奖号码，如果用户已经抽过奖则返回抽过的号码
func GetGuessMarketNumber(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	err := ctx.Bind(request)
	if err != nil {
		response.RetCode = -1
		response.RetMsg = err.Error()
		response.Message = false
		return ctx.WriteJson(response)
	}
	requestData := agent.GuessMarketNumberRequest{}
	jsonStr := _json.GetJsonString(request.RequestData)
	err = _json.Unmarshal(jsonStr, &requestData)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}

	guessMarketActivitySrv := service.NewGuessMarketActivityService()
	number, err := guessMarketActivitySrv.GetNumber(requestData.UserInfo, requestData.ActivityId)
	if err != nil {
		response.RetCode = -3
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = number

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}

// GetCurrentGuessInfo 获取当前抽奖活动信息及用户参与信息
func GetCurrentGuessInfo(ctx dotweb.Context) error {
	request := contract.NewApiRequest()
	response := contract.NewApiResponse()
	currentGuessInfo := agent.CurrentGuessInfo{}
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
	guessMarketActivitySrv := service.NewGuessMarketActivityService()
	currentGuessInfo.Activity, err = guessMarketActivitySrv.GetCurrentActivity()
	if err != nil {
		response.RetCode = -3
		response.RetMsg = "获取当前活动失败，请稍后再试"
		return ctx.WriteJson(response)
	}
	if currentGuessInfo.Activity == nil {
		response.RetCode = -4
		response.RetMsg = "当前活动不存在"
		return ctx.WriteJson(response)
	}
	usrInGuessSrv := service.NewUserInGuessMarketActivityService()
	currentGuessInfo.UserJoinRecord, err = usrInGuessSrv.GetUserInGuessMarketActivity(userInfoId, currentGuessInfo.Activity.IssueNumber)
	if err != nil {
		response.RetCode = -5
		response.RetMsg = "获取用户参与记录失败，请稍后再试"
		return ctx.WriteJson(response)
	}

	//处理用户与活动的状态
	//处理活动状态
	now := time.Now()
	if now.Before(time.Time(currentGuessInfo.Activity.EndTime)) {
		currentGuessInfo.ActivityState = _const.ActivityProgress_Begining
	} else if currentGuessInfo.Activity.IsPublish {
		currentGuessInfo.ActivityState = _const.ActivityProgress_Finish_Published
	} else {
		currentGuessInfo.ActivityState = _const.ActivityProgress_Finish_NotPublish
	}

	//处理用户参与状态
	if currentGuessInfo.UserJoinRecord == nil {
		currentGuessInfo.UserJoinState = _const.UserJoinState_NotJoin
	} else if currentGuessInfo.UserJoinRecord.IsGuess {
		currentGuessInfo.UserJoinState = _const.UserJoinState_Join_Guessed
	} else if currentGuessInfo.Activity.IsPublish {
		currentGuessInfo.UserJoinState = _const.UserJoinState_Join_Not_Guess
	} else {
		currentGuessInfo.UserJoinState = _const.UserJoinState_Join
	}

	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = currentGuessInfo

	//处理客户端缓存逻辑
	contract.SupportClientCache(request, response)
	return ctx.WriteJson(response)
}
