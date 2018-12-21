package agent

import (
	"errors"
	"git.emoney.cn/softweb/roboadvisor/config"
	"git.emoney.cn/softweb/roboadvisor/contract"
	userhome_model "git.emoney.cn/softweb/roboadvisor/protected/model/userhome"
	"git.emoney.cn/softweb/roboadvisor/util/json"
)

// GetCurrentGuessChange 获取当期猜涨跌活动情况
func GetCurrentGuessChange(userInfo userhome_model.UserInfo) (currentGuessChange *CurrentGuessChange, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userInfo
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/currentguesschange", req)
	if err != nil {
		return currentGuessChange, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return currentGuessChange, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &currentGuessChange)
	return currentGuessChange, nil
}

// GuessChangeSubmit 提交猜涨跌竞猜结果 猜测结果result 1涨 -1跌
func GuessChangeSubmit(userInfo userhome_model.UserInfo, result int) (currentGuessChange *CurrentGuessChange, err error) {
	req := contract.NewApiRequest()
	requestData := &GuessChangeSubmitRequest{
		UserInfo:   userInfo,
		Result:     result,
		ActivityId: config.CurrentConfig.Activity_GuessChange,
	}
	req.RequestData = requestData
	response, err := PostNoCache(config.CurrentConfig.WebApiHost+"/api/activity/guesschangesubmit", req)
	if err != nil {
		return currentGuessChange, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return currentGuessChange, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &currentGuessChange)
	return currentGuessChange, nil
}

// GetMyGuessChangeInfoCurrentWeek 获取本周我的竞猜记录
func GetMyGuessChangeInfoCurrentWeek(userInfo userhome_model.UserInfo) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	req := contract.NewApiRequest()
	req.RequestData = userInfo
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/myguesschangeinfocurrentweek", req)
	if err != nil {
		return myGuessChangeList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return myGuessChangeList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &myGuessChangeList)
	return myGuessChangeList, nil
}

// GetMyGuessageAward 获取我的猜涨跌奖品
func GetMyGuessageAwardList(userInfo userhome_model.UserInfo) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	req := contract.NewApiRequest()
	requestData := &MyGuessChangeInfoRequest{
		UserInfo: userInfo,
		Top:      5,
	}
	req.RequestData = requestData
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/myguessawardlist", req)
	if err != nil {
		return myGuessChangeList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return myGuessChangeList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &myGuessChangeList)
	return myGuessChangeList, nil
}

// GetMyGuessChangeInfoNewst 获取我最近的猜涨跌记录
func GetMyGuessChangeInfoNewst(userInfo userhome_model.UserInfo) (myGuessChangeList []*userhome_model.UserInGuessChangeActivity, err error) {
	req := contract.NewApiRequest()
	requestData := &MyGuessChangeInfoRequest{
		UserInfo: userInfo,
		Top:      10,
	}
	req.RequestData = requestData
	response, err := PostWithResponseCache(config.CurrentConfig.WebApiHost+"/api/activity/myguesschangeinfonewst", req)
	if err != nil {
		return myGuessChangeList, errors.New("API调用异常:" + err.Error())
	}
	if response.RetCode != 0 {
		return myGuessChangeList, errors.New(response.RetMsg)
	}
	if response.Message == nil {
		return nil, nil
	}
	jsonStr := _json.GetJsonString(response.Message)
	err = _json.Unmarshal(jsonStr, &myGuessChangeList)
	return myGuessChangeList, nil
}

//当前猜涨跌活动信息
type CurrentGuessChange struct {
	//期号 20180522
	IssueNumber string
	//我的竞猜结果  1涨 -1跌
	MyGuessResult int
	//是否已经参加了当期竞猜
	IsJoin bool
	//预计开奖时间
	//PublishTime time.Time
	//猜涨人数
	GuessUp int64
	//猜跌人数
	GuessDown int64
}

//猜涨跌提交
type GuessChangeSubmitRequest struct {
	UserInfo   userhome_model.UserInfo
	ActivityId int64
	Result     int
}

//用户参与记录查询请求
type MyGuessChangeInfoRequest struct {
	UserInfo userhome_model.UserInfo
	Top      int
}
