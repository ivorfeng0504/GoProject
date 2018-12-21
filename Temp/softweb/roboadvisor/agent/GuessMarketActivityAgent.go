package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetGuessMarketHistoryInfo 获取历史信息-包括我的抽奖记录&往期中奖号码&中奖名单
func GetGuessMarketHistoryInfo(userInfoId int) (history GuessMarketActivityHistory, err error) {
	req := contract.NewApiRequest()
	requestInfo := GuessMarketActivityHistoryRequest{
		UserInfoId:             userInfoId,
		MyGuessCount:           6,
		HistoryLuckNumberCount: 6,
		WinnerHistoryCount:     4,
	}
	req.RequestData = requestInfo
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/getguessmarkethistoryinfo", req)
	if err != nil {
		return history, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return history, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return history, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &history)
	return history, nil
}

// GetGuessMarketNumber 用户获取抽奖号码，如果用户已经抽过奖则返回抽过的号码
func GetGuessMarketNumber(userInfo userhome_model.UserInfo) (number string, err error) {
	req := contract.NewApiRequest()
	requestInfo := GuessMarketNumberRequest{
		UserInfo:   userInfo,
		ActivityId: config.CurrentConfig.Activity_GuessNumber,
	}
	req.RequestData = requestInfo
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/activity/getguessmarketnumber", req)
	if err != nil {
		return number, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return number, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return number, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &number)
	return number, nil
}

// GetCurrentGuessInfo 获取当前抽奖信息
func GetCurrentGuessInfo(userInfoId int) (currentGuessInfo CurrentGuessInfo, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userInfoId
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/getcurrentguessinfo", req)
	if err != nil {
		return currentGuessInfo, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return currentGuessInfo, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return currentGuessInfo, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &currentGuessInfo)
	return currentGuessInfo, nil
}

// 活动历史信息
type GuessMarketActivityHistory struct {
	//MyGuessCount 我的抽奖总数
	MyGuessCount int64
	// MyGuessList 我的抽奖记录 最近6条
	MyGuessList []*userhome_model.UserInGuessMarketActivity
	// HistoryLuckNumberList 往期中奖号码 最近6条
	HistoryLuckNumberList []*userhome_model.GuessMarketActivity
	// WinnerHistoryList 中奖名单 最近4条
	WinnerHistoryList []*userhome_model.UserInGuessMarketActivity
}

// 活动历史信息请求
type GuessMarketActivityHistoryRequest struct {
	// UserInfoId 用户Id
	UserInfoId int
	// MyGuessCount 我的抽奖记录 最近6条
	MyGuessCount int
	// HistoryLuckNumberCount 往期中奖号码 最近6条
	HistoryLuckNumberCount int
	// WinnerHistoryCount 中奖名单 最近4条
	WinnerHistoryCount int
}

//用户领号请求
type GuessMarketNumberRequest struct {
	UserInfo   userhome_model.UserInfo
	ActivityId int64
}

// 当前领数字活动信息
type CurrentGuessInfo struct {
	// Activity 当前活动
	Activity *userhome_model.GuessMarketActivity
	// UserJoinRecord 用户活动参与记录
	UserJoinRecord *userhome_model.UserInGuessMarketActivity
	//活动状态 详见 ActivityProgress.go
	ActivityState int
	//用户参与活动状态 详见UserJoinState.go
	UserJoinState int
}
